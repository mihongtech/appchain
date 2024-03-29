package meta

import (
	"github.com/mihongtech/appchain/protobuf"
	"github.com/mihongtech/linkchain-core/common/hexutil"
	"github.com/mihongtech/linkchain-core/common/math"
	"github.com/mihongtech/linkchain-core/common/serialize"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"

	"github.com/golang/protobuf/proto"
)

//go:generate gencodec -type Log -field-override logMarshaling -out gen_log_json.go

// Log represents a contract log event. These events are generated by the LOG opcode and
// stored/indexed by the node.
type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address node_meta.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []math.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	TxHash node_meta.TxID `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex" gencodec:"required"`
	// hash of the block in which the transaction was included
	BlockHash node_meta.BlockID `json:"blockHash"`
	// index of the log in the block
	Index uint `json:"logIndex" gencodec:"required"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"`
}

type logMarshaling struct {
	Data        hexutil.Bytes
	BlockNumber hexutil.Uint64
	TxIndex     hexutil.Uint
	Index       hexutil.Uint
}

type protoLog struct {
	Address node_meta.Address
	Topics  []math.Hash
	Data    []byte
}

type protoStorageLog struct {
	Address     node_meta.Address
	Topics      []math.Hash
	Data        []byte
	BlockNumber uint64
	TxHash      math.Hash
	TxIndex     uint
	BlockHash   math.Hash
	Index       uint
}

//Serialize/Deserialize
func (l *Log) Serialize() serialize.SerializeStream {
	topics := make([][]byte, 0)
	for i := range l.Topics {
		topic, _ := l.Topics[i].EncodeToBytes()
		topics = append(topics, topic)
	}
	bAddress, _ := l.Address.EncodeToBytes()
	log := protobuf.Log{
		Address: bAddress,
		Topics:  topics,
		Data:    proto.NewBuffer(l.Data).Bytes(),
	}
	return &log
}

func (l *Log) Deserialize(s serialize.SerializeStream) error {
	data := s.(*protobuf.Log)
	if err := l.Address.DecodeFromBytes(data.Address); err != nil {
		return err
	}
	l.Data = l.Data[:0]
	copy(l.Data, data.Data)

	l.Topics = l.Topics[:0]
	for i := range data.Topics {
		topic := math.Hash{}
		if err := topic.DecodeFromBytes(data.Topics[i]); err != nil {
			return err
		}
		l.Topics = append(l.Topics, topic)
	}

	return nil
}

// LogForStorage is a wrapper around a Log that flattens and parses the entire content of
// a log including non-consensus fields.
type LogForStorage Log

//Serialize/Deserialize
func (l *LogForStorage) Serialize() serialize.SerializeStream {
	topics := make([][]byte, 0)

	for i := range l.Topics {
		topic, _ := l.Topics[i].EncodeToBytes()
		topics = append(topics, topic)
	}
	bAddress, _ := l.Address.EncodeToBytes()
	bBlockHash, _ := l.BlockHash.EncodeToBytes()
	bTxHash, _ := l.TxHash.EncodeToBytes()
	log := protobuf.LogForStorage{
		Address:     bAddress,
		Topics:      topics,
		Data:        proto.NewBuffer(l.Data).Bytes(),
		BlockNumber: proto.Uint64(l.BlockNumber),
		BlockHash:   bBlockHash,
		Index:       proto.Uint32(uint32(l.Index)),
		TxHash:      bTxHash,
		TxIndex:     proto.Uint32(uint32(l.TxIndex)),
	}
	return &log
}

func (l *LogForStorage) Deserialize(s serialize.SerializeStream) error {
	data := s.(*protobuf.LogForStorage)

	l.Index = uint(*data.Index)
	l.TxIndex = uint(*data.TxIndex)
	l.BlockNumber = *data.BlockNumber
	if err := l.TxHash.DecodeFromBytes(data.TxHash); err != nil {
		return err
	}
	if err := l.BlockHash.DecodeFromBytes(data.BlockHash); err != nil {
		return err
	}

	if err := l.Address.DecodeFromBytes(data.Address); err != nil {
		return err
	}
	l.Data = l.Data[:0]
	copy(l.Data, data.Data)

	l.Topics = l.Topics[:0]
	for i := range data.Topics {
		topic := math.Hash{}
		if err := topic.DecodeFromBytes(data.Topics[i]); err != nil {
			return err
		}
		l.Topics = append(l.Topics, topic)
	}

	return nil
}

func (l *LogForStorage) ConvertToLog() *Log {
	log := Log{}
	log.BlockHash, log.TxHash, log.BlockNumber, log.TxIndex, log.Index, log.Topics, log.Data, log.Address = l.BlockHash, l.TxHash, l.BlockNumber, l.TxIndex, l.Index, l.Topics, l.Data, l.Address
	return &log
}
