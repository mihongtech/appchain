package rpcserver

import (
	"fmt"
	"github.com/mihongtech/linkchain-core/common/util/log"
	"reflect"

	"github.com/mihongtech/appchain/rpc/rpcobject"
	p2p_node "github.com/mihongtech/linkchain-core/node/net/p2p/discover"
)

func selfPeer(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	server := GetP2PAPI(s)
	self := server.Self()

	return self.String(), nil
}

func addPeer(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c, ok := cmd.(*rpcobject.PeerCmd)
	if !ok {
		fmt.Println("Type error:", reflect.TypeOf(cmd))
		return nil, nil
	}

	server := GetP2PAPI(s)
	node, err := p2p_node.ParseNode(c.Peer)
	if err != nil {
		log.Error("parse node failes", "url", c.Peer)
		return nil, err
	}
	server.AddPeer(node)

	return nil, nil
}

func listPeer(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	server := GetP2PAPI(s)
	peers := server.Peers()

	ps := make([]string, 0)
	for _, peer := range peers {
		ps = append(ps, peer.String())
	}

	return ps, nil
}

func removePeer(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c, ok := cmd.(*rpcobject.PeerCmd)
	if !ok {
		fmt.Println("Type error:", reflect.TypeOf(cmd))
		return nil, nil
	}

	server := GetP2PAPI(s)
	node, err := p2p_node.ParseNode(c.Peer)
	if err != nil {
		log.Error("par"+
			""+
			"se node failes", "url", c.Peer)
		return nil, err
	}
	server.RemovePeer(node)

	return nil, nil
}
