package discover

import (
	"bytes"
	"container/list"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/linkchain/common/btcec"
	"github.com/linkchain/common/math"
	"github.com/linkchain/common/serialize"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/p2p/netutil"
	"github.com/linkchain/protobuf"

	"github.com/golang/protobuf/proto"
)

const Version = 4

// Errors
var (
	errPacketTooSmall   = errors.New("too small")
	errBadHash          = errors.New("bad hash")
	errExpired          = errors.New("expired")
	errUnsolicitedReply = errors.New("unsolicited reply")
	errUnknownNode      = errors.New("unknown node")
	errTimeout          = errors.New("RPC timeout")
	errClockWarp        = errors.New("reply deadline too far in the future")
	errClosed           = errors.New("socket closed")
)

// Timeouts
const (
	respTimeout = 500 * time.Millisecond
	sendTimeout = 500 * time.Millisecond
	expiration  = 20 * time.Second

	ntpFailureThreshold = 32               // Continuous timeouts after which to check NTP
	ntpWarningCooldown  = 10 * time.Minute // Minimum amount of time to pass before repeating NTP warning
	driftThreshold      = 10 * time.Second // Allowed clock drift before warning user
)

// RPC packet types
const (
	pingPacket = iota + 1 // zero is 'reserved'
	pongPacket
	findnodePacket
	neighborsPacket
)

// RPC request structures
type (
	ping struct {
		Version    uint
		From, To   rpcEndpoint
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []byte `rlp:"tail"`
	}

	// pong is the reply to ping.
	pong struct {
		// This field should mirror the UDP envelope address
		// of the ping packet, which provides a way to discover the
		// the external address (after NAT).
		To rpcEndpoint

		ReplyTok   []byte // This contains the hash of the ping packet.
		Expiration uint64 // Absolute timestamp at which the packet becomes invalid.
		// Ignore additional fields (for forward compatibility).
		Rest []byte `rlp:"tail"`
	}

	// findnode is a query for nodes close to the given target.
	findnode struct {
		Target     NodeID // doesn't need to be an actual public key
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []byte `rlp:"tail"`
	}

	// reply to findnode
	neighbors struct {
		Nodes      []rpcNode
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []byte `rlp:"tail"`
	}

	rpcNode struct {
		IP  net.IP // len 4 for IPv4 or 16 for IPv6
		UDP uint16 // for discovery protocol
		TCP uint16 // for RLPx protocol
		ID  NodeID
	}

	rpcEndpoint struct {
		IP  net.IP // len 4 for IPv4 or 16 for IPv6
		UDP uint16 // for discovery protocol
		TCP uint16 // for RLPx protocol
	}
)

func (n *rpcNode) Serialize() serialize.SerializeStream {
	udp := uint32(n.UDP)
	tcp := uint32(n.TCP)
	node := protobuf.Node{
		Ip:  n.IP,
		Udp: &udp,
		Tcp: &tcp,
		Id:  n.ID[:],
	}

	return &node
}

func (n *rpcNode) Deserialize(s serialize.SerializeStream) error {
	node := *s.(*protobuf.Node)

	n.ID, _ = BytesID(node.Id)
	n.IP = node.Ip
	n.TCP = uint16(*node.Tcp)
	n.UDP = uint16(*node.Udp)

	return nil
}

func (e *rpcEndpoint) Serialize() serialize.SerializeStream {
	udp := uint32(e.UDP)
	tcp := uint32(e.TCP)
	ip := []byte{}
	if e.IP != nil {
		ip = e.IP
	}
	endpoint := protobuf.Endpoint{
		Ip:  ip,
		Udp: &udp,
		Tcp: &tcp,
	}

	return &endpoint
}

func (e *rpcEndpoint) Deserialize(s serialize.SerializeStream) error {
	node := *s.(*protobuf.Endpoint)

	e.IP = node.Ip
	e.TCP = uint16(*node.Tcp)
	e.UDP = uint16(*node.Udp)

	return nil
}

func (p *ping) Serialize() serialize.SerializeStream {

	version := uint64(p.Version)
	from := p.From.Serialize()
	to := p.To.Serialize()
	pingMsg := protobuf.Ping{
		Version:    &version,
		From:       from.(*protobuf.Endpoint),
		To:         to.(*protobuf.Endpoint),
		Expiration: &p.Expiration,
		Rest:       p.Rest,
	}

	return &pingMsg
}

func (p *ping) Deserialize(s serialize.SerializeStream) error {
	pingMsg := *s.(*protobuf.Ping)

	p.Version = uint(*pingMsg.Version)
	p.From.Deserialize(pingMsg.From)
	p.To.Deserialize(pingMsg.To)
	p.Expiration = *pingMsg.Expiration
	p.Rest = pingMsg.Rest

	return nil
}

func (p *pong) Serialize() serialize.SerializeStream {

	to := p.To.Serialize()
	pongMsg := protobuf.Pong{
		To:         to.(*protobuf.Endpoint),
		ReplyTok:   p.ReplyTok,
		Expiration: &p.Expiration,
		Rest:       p.Rest,
	}

	return &pongMsg
}

func (p *pong) Deserialize(s serialize.SerializeStream) error {
	pongMsg := *s.(*protobuf.Pong)

	p.To.Deserialize(pongMsg.To)
	p.ReplyTok = pongMsg.ReplyTok
	p.Expiration = *pongMsg.Expiration
	p.Rest = pongMsg.Rest

	return nil
}

func (f *findnode) Serialize() serialize.SerializeStream {

	findMsg := protobuf.Findnode{
		Target:     f.Target.Bytes(),
		Expiration: &f.Expiration,
		Rest:       f.Rest,
	}

	return &findMsg
}

func (f *findnode) Deserialize(s serialize.SerializeStream) error {
	findMsg := *s.(*protobuf.Findnode)
	var err error
	f.Target, err = BytesID(findMsg.Target)
	if err != nil {
		return err
	}
	f.Expiration = *findMsg.Expiration
	f.Rest = findMsg.Rest

	return nil
}

func (n *neighbors) Serialize() serialize.SerializeStream {
	nodes := make([]*protobuf.Node, 0, 0)
	for i := range n.Nodes {
		nodes = append(nodes, n.Nodes[i].Serialize().(*protobuf.Node))
	}
	neighborsMsg := protobuf.Neighbors{
		Nodes:      nodes,
		Expiration: &n.Expiration,
		Rest:       n.Rest,
	}

	return &neighborsMsg
}

func (n *neighbors) Deserialize(s serialize.SerializeStream) error {
	neighborsMsg := *s.(*protobuf.Neighbors)
	n.Nodes = make([]rpcNode, 0, 0)
	n.Expiration = *neighborsMsg.Expiration
	n.Rest = neighborsMsg.Rest
	for i := range neighborsMsg.Nodes {
		node := new(rpcNode)
		node.Deserialize(neighborsMsg.Nodes[i])
		n.Nodes = append(n.Nodes, *node)
	}

	return nil
}

func makeEndpoint(addr *net.UDPAddr, tcpPort uint16) rpcEndpoint {
	ip := addr.IP.To4()
	if ip == nil {
		ip = addr.IP.To16()
	}
	return rpcEndpoint{IP: ip, UDP: uint16(addr.Port), TCP: tcpPort}
}

func (t *udp) nodeFromRPC(sender *net.UDPAddr, rn rpcNode) (*Node, error) {
	if rn.UDP <= 1024 {
		return nil, errors.New("low port")
	}
	if err := netutil.CheckRelayIP(sender.IP, rn.IP); err != nil {
		return nil, err
	}
	if t.netrestrict != nil && !t.netrestrict.Contains(rn.IP) {
		return nil, errors.New("not contained in netrestrict whitelist")
	}
	n := NewNode(rn.ID, rn.IP, rn.UDP, rn.TCP)
	err := n.validateComplete()
	return n, err
}

func nodeToRPC(n *Node) rpcNode {
	return rpcNode{ID: n.ID, IP: n.IP, UDP: n.UDP, TCP: n.TCP}
}

type packet interface {
	handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error
	name() string
	serialize.ISerialize
}

type conn interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
	Close() error
	LocalAddr() net.Addr
}

// udp implements the RPC protocol.
type udp struct {
	conn        conn
	netrestrict *netutil.Netlist
	priv        *ecdsa.PrivateKey
	ourEndpoint rpcEndpoint

	addpending chan *pending
	gotreply   chan reply

	closing chan struct{}

	*Table
}

// pending represents a pending reply.
//
// some implementations of the protocol wish to send more than one
// reply packet to findnode. in general, any neighbors packet cannot
// be matched up with a specific findnode packet.
//
// our implementation handles this by storing a callback function for
// each pending reply. incoming packets from a node are dispatched
// to all the callback functions for that node.
type pending struct {
	// these fields must match in the reply.
	from  NodeID
	ptype byte

	// time when the request must complete
	deadline time.Time

	// callback is called when a matching reply arrives. if it returns
	// true, the callback is removed from the pending reply queue.
	// if it returns false, the reply is considered incomplete and
	// the callback will be invoked again for the next matching reply.
	callback func(resp interface{}) (done bool)

	// errc receives nil when the callback indicates completion or an
	// error if no further reply is received within the timeout.
	errc chan<- error
}

type reply struct {
	from  NodeID
	ptype byte
	data  interface{}
	// loop indicates whether there was
	// a matching request by sending on this channel.
	matched chan<- bool
}

// ReadPacket is sent to the unhandled channel when it could not be processed
type ReadPacket struct {
	Data []byte
	Addr *net.UDPAddr
}

// Config holds Table-related settings.
type Config struct {
	// These settings are required and configure the UDP listener:
	PrivateKey *ecdsa.PrivateKey

	// These settings are optional:
	AnnounceAddr *net.UDPAddr      // local address announced in the DHT
	NodeDBPath   string            // if set, the node database is stored at this filesystem location
	NetRestrict  *netutil.Netlist  // network whitelist
	Bootnodes    []*Node           // list of bootstrap nodes
	Unhandled    chan<- ReadPacket // unhandled packets are sent on this channel
}

// ListenUDP returns a new table that listens for UDP packets on laddr.
func ListenUDP(c conn, cfg Config) (*Table, error) {
	tab, _, err := newUDP(c, cfg)
	if err != nil {
		return nil, err
	}
	log.Info("UDP listener up", "self", tab.self)
	return tab, nil
}

func newUDP(c conn, cfg Config) (*Table, *udp, error) {
	udp := &udp{
		conn:        c,
		priv:        cfg.PrivateKey,
		netrestrict: cfg.NetRestrict,
		closing:     make(chan struct{}),
		gotreply:    make(chan reply),
		addpending:  make(chan *pending),
	}
	realaddr := c.LocalAddr().(*net.UDPAddr)
	if cfg.AnnounceAddr != nil {
		realaddr = cfg.AnnounceAddr
	}
	// TODO: separate TCP port
	udp.ourEndpoint = makeEndpoint(realaddr, uint16(realaddr.Port))
	tab, err := newTable(udp, PubkeyID(&cfg.PrivateKey.PublicKey), realaddr, cfg.NodeDBPath, cfg.Bootnodes)
	if err != nil {
		return nil, nil, err
	}
	udp.Table = tab

	go udp.loop()
	go udp.readLoop(cfg.Unhandled)
	return udp.Table, udp, nil
}

func (t *udp) close() {
	close(t.closing)
	t.conn.Close()
	// TODO: wait for the loops to end.
}

// ping sends a ping message to the given node and waits for a reply.
func (t *udp) ping(toid NodeID, toaddr *net.UDPAddr) error {
	req := &ping{
		Version:    Version,
		From:       t.ourEndpoint,
		To:         makeEndpoint(toaddr, 0), // TODO: maybe use known TCP port from DB
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	}

	packet, hash, err := encodePacket(t.priv, pingPacket, req.Serialize())
	if err != nil {
		return err
	}
	errc := t.pending(toid, pongPacket, func(p interface{}) bool {
		log.Trace("callback", "hash", hash, "p.(*pong).ReplyTok", p.(*pong).ReplyTok, "result", bytes.Equal(p.(*pong).ReplyTok, hash))
		return bytes.Equal(p.(*pong).ReplyTok, hash)
	})
	t.write(toaddr, req.name(), packet)
	return <-errc
}

func (t *udp) waitping(from NodeID) error {
	return <-t.pending(from, pingPacket, func(interface{}) bool { return true })
}

// findnode sends a findnode request to the given node and waits until
// the node has sent up to k neighbors.
func (t *udp) findnode(toid NodeID, toaddr *net.UDPAddr, target NodeID) ([]*Node, error) {
	nodes := make([]*Node, 0, bucketSize)
	nreceived := 0
	errc := t.pending(toid, neighborsPacket, func(r interface{}) bool {
		reply := r.(*neighbors)
		log.Trace("reply findnodes", "neighbors", reply.Nodes, "reply.Rest", reply.Rest)
		for _, rn := range reply.Nodes {
			nreceived++
			n, err := t.nodeFromRPC(toaddr, rn)
			if err != nil {
				log.Trace("Invalid neighbor node received", "ip", rn.IP, "addr", toaddr, "err", err)
				continue
			}
			nodes = append(nodes, n)
		}
		return nreceived >= bucketSize
	})
	t.send(toaddr, findnodePacket, &findnode{
		Target:     target,
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	})
	err := <-errc
	return nodes, err
}

// pending adds a reply callback to the pending reply queue.
// see the documentation of type pending for a detailed explanation.
func (t *udp) pending(id NodeID, ptype byte, callback func(interface{}) bool) <-chan error {
	ch := make(chan error, 1)
	p := &pending{from: id, ptype: ptype, callback: callback, errc: ch}
	select {
	case t.addpending <- p:
		// loop will handle it
	case <-t.closing:
		ch <- errClosed
	}
	return ch
}

func (t *udp) handleReply(from NodeID, ptype byte, req packet) bool {
	matched := make(chan bool, 1)
	select {
	case t.gotreply <- reply{from, ptype, req, matched}:
		// loop will handle it
		return <-matched
	case <-t.closing:
		return false
	}
}

// loop runs in its own goroutine. it keeps track of
// the refresh timer and the pending reply queue.
func (t *udp) loop() {
	var (
		plist        = list.New()
		timeout      = time.NewTimer(0)
		nextTimeout  *pending // head of plist when timeout was last reset
		contTimeouts = 0      // number of continuous timeouts to do NTP checks
		ntpWarnTime  = time.Unix(0, 0)
	)
	<-timeout.C // ignore first timeout
	defer timeout.Stop()

	resetTimeout := func() {
		if plist.Front() == nil || nextTimeout == plist.Front().Value {
			return
		}
		// Start the timer so it fires when the next pending reply has expired.
		now := time.Now()
		for el := plist.Front(); el != nil; el = el.Next() {
			nextTimeout = el.Value.(*pending)
			if dist := nextTimeout.deadline.Sub(now); dist < 2*respTimeout {
				timeout.Reset(dist)
				return
			}
			// Remove pending replies whose deadline is too far in the
			// future. These can occur if the system clock jumped
			// backwards after the deadline was assigned.
			nextTimeout.errc <- errClockWarp
			plist.Remove(el)
		}
		nextTimeout = nil
		timeout.Stop()
	}

	for {
		resetTimeout()

		select {
		case <-t.closing:
			for el := plist.Front(); el != nil; el = el.Next() {
				el.Value.(*pending).errc <- errClosed
			}
			return

		case p := <-t.addpending:
			p.deadline = time.Now().Add(respTimeout)
			plist.PushBack(p)

		case r := <-t.gotreply:
			log.Trace("get reply", "r", r)
			var matched bool
			for el := plist.Front(); el != nil; el = el.Next() {
				p := el.Value.(*pending)
				if p.from == r.from && p.ptype == r.ptype {
					matched = true
					// Remove the matcher if its callback indicates
					// that all replies have been received. This is
					// required for packet types that expect multiple
					// reply packets.
					if p.callback(r.data) {
						p.errc <- nil
						plist.Remove(el)
					}
					// Reset the continuous timeout counter (time drift detection)
					contTimeouts = 0
				}
			}
			r.matched <- matched
			log.Trace("get reply", "matched", matched)
		case now := <-timeout.C:
			nextTimeout = nil

			// Notify and remove callbacks whose deadline is in the past.
			for el := plist.Front(); el != nil; el = el.Next() {
				p := el.Value.(*pending)
				if now.After(p.deadline) || now.Equal(p.deadline) {
					log.Trace("timeout", "now", now, "p.deadline", p.deadline)
					p.errc <- errTimeout
					plist.Remove(el)
					contTimeouts++
				}
			}
			// If we've accumulated too many timeouts, do an NTP time sync check
			if contTimeouts > ntpFailureThreshold {
				if time.Since(ntpWarnTime) >= ntpWarningCooldown {
					ntpWarnTime = time.Now()
					go checkClockDrift()
				}
				contTimeouts = 0
			}
		}
	}
}

const (
	macSize  = 256 / 8
	sigSize  = 520 / 8
	headSize = macSize + sigSize // space of packet frame data
)

var (
	headSpace = make([]byte, headSize)

	// Neighbors replies are sent across multiple packets to
	// stay below the 1280 byte limit. We compute the maximum number
	// of entries by stuffing a packet until it grows too large.
	maxNeighbors int
)

func init() {
	p := neighbors{Expiration: ^uint64(0)}
	maxSizeNode := rpcNode{IP: make(net.IP, 16), UDP: ^uint16(0), TCP: ^uint16(0)}
	for n := 0; ; n++ {
		p.Nodes = append(p.Nodes, maxSizeNode)
		data, err := proto.Marshal(p.Serialize())
		if err != nil {
			// If this ever happens, it will be caught by the unit tests.
			panic("cannot encode: " + err.Error())
		}
		size := len(data)
		if headSize+size+1 >= 1280 {
			maxNeighbors = n
			break
		}
	}
}

func (t *udp) send(toaddr *net.UDPAddr, ptype byte, req packet) ([]byte, error) {
	packet, hash, err := encodePacket(t.priv, ptype, req.Serialize())
	if err != nil {
		return hash, err
	}
	return hash, t.write(toaddr, req.name(), packet)
}

func (t *udp) write(toaddr *net.UDPAddr, what string, packet []byte) error {
	_, err := t.conn.WriteToUDP(packet, toaddr)
	log.Trace(">> "+what, "addr", toaddr, "ID", t.Self().ID, "err", err)
	return err
}

func encodePacket(priv *ecdsa.PrivateKey, ptype byte, req proto.Message) (packet, hash []byte, err error) {
	b := new(bytes.Buffer)
	b.Write(headSpace)
	b.WriteByte(ptype)
	data, err := proto.Marshal(req)
	if err != nil {
		log.Error("Can't encode discv4 packet", "err", err)
		return nil, nil, err
	}
	b.Write(data)
	packet = b.Bytes()
	signData, err := btcec.SignCompact(btcec.S256(), (*btcec.PrivateKey)(priv), math.HashB(packet[headSize:]), true)
	if err != nil {
		log.Error("Can't sign discv4 packet", "err", err)
		return nil, nil, err
	}

	copy(packet[macSize:], signData)
	// add the hash to the front. Note: this doesn't protect the
	// packet in any way. Our public key will be part of this hash in
	// The future.
	hash = math.HashB(packet[macSize:])
	copy(packet, hash)
	return packet, hash, nil
}

// readLoop runs in its own goroutine. it handles incoming UDP packets.
func (t *udp) readLoop(unhandled chan<- ReadPacket) {
	defer t.conn.Close()
	if unhandled != nil {
		defer close(unhandled)
	}
	// Discovery packets are defined to be no larger than 1280 bytes.
	// Packets larger than this size will be cut at the end and treated
	// as invalid because their hash won't match.
	buf := make([]byte, 1280)
	for {
		nbytes, from, err := t.conn.ReadFromUDP(buf)
		if netutil.IsTemporaryError(err) {
			// Ignore temporary read errors.
			log.Debug("Temporary UDP read error", "err", err)
			continue
		} else if err != nil {
			// Shut down the loop for permament errors.
			log.Debug("UDP read error", "err", err)
			return
		}
		if t.handlePacket(from, buf[:nbytes]) != nil && unhandled != nil {
			select {
			case unhandled <- ReadPacket{buf[:nbytes], from}:
			default:
			}
		}
	}
}

func (t *udp) handlePacket(from *net.UDPAddr, buf []byte) error {
	packet, fromID, hash, err := decodePacket(buf)
	if err != nil {
		log.Debug("Bad discv4 packet", "addr", from, "err", err)
		return err
	}
	err = packet.handle(t, from, fromID, hash)
	log.Trace("<< "+packet.name(), "addr", from, "formID", fromID, "err", err)
	return err
}

func decodePacket(buf []byte) (packet, NodeID, []byte, error) {
	if len(buf) < headSize+1 {
		return nil, NodeID{}, nil, errPacketTooSmall
	}
	hash, sig, sigdata := buf[:macSize], buf[macSize:headSize], buf[headSize:]
	shouldhash := math.HashB(buf[macSize:])
	if !bytes.Equal(hash, shouldhash) {
		return nil, NodeID{}, nil, errBadHash
	}
	fromID, err := recoverNodeID(math.HashB(buf[headSize:]), sig)
	if err != nil {
		return nil, NodeID{}, hash, err
	}
	var req packet
	var msg proto.Message
	switch ptype := sigdata[0]; ptype {
	case pingPacket:
		req = new(ping)
		msg = new(protobuf.Ping)
	case pongPacket:
		req = new(pong)
		msg = new(protobuf.Pong)
	case findnodePacket:
		req = new(findnode)
		msg = new(protobuf.Findnode)
	case neighborsPacket:
		req = new(neighbors)
		msg = new(protobuf.Neighbors)
	default:
		return nil, fromID, hash, fmt.Errorf("unknown type: %d", ptype)
	}

	err = proto.Unmarshal(sigdata[1:], msg)
	req.Deserialize(msg)
	return req, fromID, hash, err
}

func (req *ping) handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error {
	if expired(req.Expiration) {
		return errExpired
	}
	t.send(from, pongPacket, &pong{
		To:         makeEndpoint(from, req.From.TCP),
		ReplyTok:   mac,
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	})
	if !t.handleReply(fromID, pingPacket, req) {
		// Note: we're ignoring the provided IP address right now
		go t.bond(true, fromID, from, req.From.TCP)
	}
	return nil
}

func (req *ping) name() string   { return "PING/v4" }
func (req *ping) String() string { return "PING/v4" }

func (req *pong) handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error {
	if expired(req.Expiration) {
		return errExpired
	}
	if !t.handleReply(fromID, pongPacket, req) {
		return errUnsolicitedReply
	}
	return nil
}

func (req *pong) name() string { return "PONG/v4" }

func (req *pong) String() string { return "PONG/v4" }

func (req *findnode) handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error {
	if expired(req.Expiration) {
		return errExpired
	}
	if !t.db.hasBond(fromID) {
		// No bond exists, we don't process the packet. This prevents
		// an attack vector where the discovery protocol could be used
		// to amplify traffic in a DDOS attack. A malicious actor
		// would send a findnode request with the IP address and UDP
		// port of the target as the source address. The recipient of
		// the findnode packet would then send a neighbors packet
		// (which is a much bigger packet than findnode) to the victim.
		return errUnknownNode
	}
	target := math.HashH(req.Target[:])
	t.mutex.Lock()
	closest := t.closest(target, bucketSize).entries
	t.mutex.Unlock()

	p := neighbors{Expiration: uint64(time.Now().Add(expiration).Unix())}
	var sent bool
	// Send neighbors in chunks with at most maxNeighbors per packet
	// to stay below the 1280 byte limit.
	for _, n := range closest {
		if netutil.CheckRelayIP(from.IP, n.IP) == nil {
			p.Nodes = append(p.Nodes, nodeToRPC(n))
		}
		if len(p.Nodes) == maxNeighbors {
			t.send(from, neighborsPacket, &p)
			p.Nodes = p.Nodes[:0]
			sent = true
		}
	}
	if len(p.Nodes) > 0 || !sent {
		t.send(from, neighborsPacket, &p)
	}
	return nil
}

func (req *findnode) name() string   { return "FINDNODE/v4" }
func (req *findnode) String() string { return "FINDNODE/v4" }

func (req *neighbors) handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error {
	if expired(req.Expiration) {
		return errExpired
	}
	if !t.handleReply(fromID, neighborsPacket, req) {
		return errUnsolicitedReply
	}
	return nil
}

func (req *neighbors) name() string   { return "NEIGHBORS/v4" }
func (req *neighbors) String() string { return "NEIGHBORS/v4" }

func expired(ts uint64) bool {
	return time.Unix(int64(ts), 0).Before(time.Now())
}
