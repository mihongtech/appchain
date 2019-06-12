// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/contract.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type BlockHeaderData struct {
	ReceiptHash          []byte   `protobuf:"bytes,1,req,name=receiptHash" json:"receiptHash,omitempty"`
	GasLimit             *uint64  `protobuf:"varint,2,req,name=gasLimit" json:"gasLimit,omitempty"`
	GasUsed              *uint64  `protobuf:"varint,3,req,name=gasUsed" json:"gasUsed,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockHeaderData) Reset()         { *m = BlockHeaderData{} }
func (m *BlockHeaderData) String() string { return proto.CompactTextString(m) }
func (*BlockHeaderData) ProtoMessage()    {}
func (*BlockHeaderData) Descriptor() ([]byte, []int) {
	return fileDescriptor_0afd3d30283bce23, []int{0}
}

func (m *BlockHeaderData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockHeaderData.Unmarshal(m, b)
}
func (m *BlockHeaderData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockHeaderData.Marshal(b, m, deterministic)
}
func (m *BlockHeaderData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockHeaderData.Merge(m, src)
}
func (m *BlockHeaderData) XXX_Size() int {
	return xxx_messageInfo_BlockHeaderData.Size(m)
}
func (m *BlockHeaderData) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockHeaderData.DiscardUnknown(m)
}

var xxx_messageInfo_BlockHeaderData proto.InternalMessageInfo

func (m *BlockHeaderData) GetReceiptHash() []byte {
	if m != nil {
		return m.ReceiptHash
	}
	return nil
}

func (m *BlockHeaderData) GetGasLimit() uint64 {
	if m != nil && m.GasLimit != nil {
		return *m.GasLimit
	}
	return 0
}

func (m *BlockHeaderData) GetGasUsed() uint64 {
	if m != nil && m.GasUsed != nil {
		return *m.GasUsed
	}
	return 0
}

type TxData struct {
	Price                *uint64  `protobuf:"varint,1,req,name=price" json:"price,omitempty"`
	GasLimit             *uint64  `protobuf:"varint,2,req,name=gasLimit" json:"gasLimit,omitempty"`
	Payload              []byte   `protobuf:"bytes,3,req,name=payload" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxData) Reset()         { *m = TxData{} }
func (m *TxData) String() string { return proto.CompactTextString(m) }
func (*TxData) ProtoMessage()    {}
func (*TxData) Descriptor() ([]byte, []int) {
	return fileDescriptor_0afd3d30283bce23, []int{1}
}

func (m *TxData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxData.Unmarshal(m, b)
}
func (m *TxData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxData.Marshal(b, m, deterministic)
}
func (m *TxData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxData.Merge(m, src)
}
func (m *TxData) XXX_Size() int {
	return xxx_messageInfo_TxData.Size(m)
}
func (m *TxData) XXX_DiscardUnknown() {
	xxx_messageInfo_TxData.DiscardUnknown(m)
}

var xxx_messageInfo_TxData proto.InternalMessageInfo

func (m *TxData) GetPrice() uint64 {
	if m != nil && m.Price != nil {
		return *m.Price
	}
	return 0
}

func (m *TxData) GetGasLimit() uint64 {
	if m != nil && m.GasLimit != nil {
		return *m.GasLimit
	}
	return 0
}

func (m *TxData) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type Receipt struct {
	PostStateOrStatus    []byte   `protobuf:"bytes,1,req,name=postStateOrStatus" json:"postStateOrStatus,omitempty"`
	CumulativeGasUsed    *uint64  `protobuf:"varint,2,req,name=cumulativeGasUsed" json:"cumulativeGasUsed,omitempty"`
	Bloom                []byte   `protobuf:"bytes,3,req,name=bloom" json:"bloom,omitempty"`
	Logs                 []*Log   `protobuf:"bytes,4,rep,name=logs" json:"logs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Receipt) Reset()         { *m = Receipt{} }
func (m *Receipt) String() string { return proto.CompactTextString(m) }
func (*Receipt) ProtoMessage()    {}
func (*Receipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_0afd3d30283bce23, []int{2}
}

func (m *Receipt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Receipt.Unmarshal(m, b)
}
func (m *Receipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Receipt.Marshal(b, m, deterministic)
}
func (m *Receipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Receipt.Merge(m, src)
}
func (m *Receipt) XXX_Size() int {
	return xxx_messageInfo_Receipt.Size(m)
}
func (m *Receipt) XXX_DiscardUnknown() {
	xxx_messageInfo_Receipt.DiscardUnknown(m)
}

var xxx_messageInfo_Receipt proto.InternalMessageInfo

func (m *Receipt) GetPostStateOrStatus() []byte {
	if m != nil {
		return m.PostStateOrStatus
	}
	return nil
}

func (m *Receipt) GetCumulativeGasUsed() uint64 {
	if m != nil && m.CumulativeGasUsed != nil {
		return *m.CumulativeGasUsed
	}
	return 0
}

func (m *Receipt) GetBloom() []byte {
	if m != nil {
		return m.Bloom
	}
	return nil
}

func (m *Receipt) GetLogs() []*Log {
	if m != nil {
		return m.Logs
	}
	return nil
}

type Log struct {
	Address              []byte   `protobuf:"bytes,1,req,name=address" json:"address,omitempty"`
	Topics               [][]byte `protobuf:"bytes,2,rep,name=topics" json:"topics,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,req,name=data" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Log) Reset()         { *m = Log{} }
func (m *Log) String() string { return proto.CompactTextString(m) }
func (*Log) ProtoMessage()    {}
func (*Log) Descriptor() ([]byte, []int) {
	return fileDescriptor_0afd3d30283bce23, []int{3}
}

func (m *Log) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Log.Unmarshal(m, b)
}
func (m *Log) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Log.Marshal(b, m, deterministic)
}
func (m *Log) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Log.Merge(m, src)
}
func (m *Log) XXX_Size() int {
	return xxx_messageInfo_Log.Size(m)
}
func (m *Log) XXX_DiscardUnknown() {
	xxx_messageInfo_Log.DiscardUnknown(m)
}

var xxx_messageInfo_Log proto.InternalMessageInfo

func (m *Log) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Log) GetTopics() [][]byte {
	if m != nil {
		return m.Topics
	}
	return nil
}

func (m *Log) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReceiptForStorage struct {
	PostStateOrStatus    []byte           `protobuf:"bytes,1,req,name=postStateOrStatus" json:"postStateOrStatus,omitempty"`
	CumulativeGasUsed    *uint64          `protobuf:"varint,2,req,name=cumulativeGasUsed" json:"cumulativeGasUsed,omitempty"`
	Bloom                []byte           `protobuf:"bytes,3,req,name=bloom" json:"bloom,omitempty"`
	Logs                 []*LogForStorage `protobuf:"bytes,4,rep,name=logs" json:"logs,omitempty"`
	TxHash               []byte           `protobuf:"bytes,5,req,name=txHash" json:"txHash,omitempty"`
	GasUsed              *uint64          `protobuf:"varint,6,req,name=gasUsed" json:"gasUsed,omitempty"`
	ContractAddress      []byte           `protobuf:"bytes,7,req,name=ContractAddress" json:"ContractAddress,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ReceiptForStorage) Reset()         { *m = ReceiptForStorage{} }
func (m *ReceiptForStorage) String() string { return proto.CompactTextString(m) }
func (*ReceiptForStorage) ProtoMessage()    {}
func (*ReceiptForStorage) Descriptor() ([]byte, []int) {
	return fileDescriptor_0afd3d30283bce23, []int{4}
}

func (m *ReceiptForStorage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReceiptForStorage.Unmarshal(m, b)
}
func (m *ReceiptForStorage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReceiptForStorage.Marshal(b, m, deterministic)
}
func (m *ReceiptForStorage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReceiptForStorage.Merge(m, src)
}
func (m *ReceiptForStorage) XXX_Size() int {
	return xxx_messageInfo_ReceiptForStorage.Size(m)
}
func (m *ReceiptForStorage) XXX_DiscardUnknown() {
	xxx_messageInfo_ReceiptForStorage.DiscardUnknown(m)
}

var xxx_messageInfo_ReceiptForStorage proto.InternalMessageInfo

func (m *ReceiptForStorage) GetPostStateOrStatus() []byte {
	if m != nil {
		return m.PostStateOrStatus
	}
	return nil
}

func (m *ReceiptForStorage) GetCumulativeGasUsed() uint64 {
	if m != nil && m.CumulativeGasUsed != nil {
		return *m.CumulativeGasUsed
	}
	return 0
}

func (m *ReceiptForStorage) GetBloom() []byte {
	if m != nil {
		return m.Bloom
	}
	return nil
}

func (m *ReceiptForStorage) GetLogs() []*LogForStorage {
	if m != nil {
		return m.Logs
	}
	return nil
}

func (m *ReceiptForStorage) GetTxHash() []byte {
	if m != nil {
		return m.TxHash
	}
	return nil
}

func (m *ReceiptForStorage) GetGasUsed() uint64 {
	if m != nil && m.GasUsed != nil {
		return *m.GasUsed
	}
	return 0
}

func (m *ReceiptForStorage) GetContractAddress() []byte {
	if m != nil {
		return m.ContractAddress
	}
	return nil
}

type ReceiptForStorages struct {
	Receipts             []*ReceiptForStorage `protobuf:"bytes,1,rep,name=receipts" json:"receipts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ReceiptForStorages) Reset()         { *m = ReceiptForStorages{} }
func (m *ReceiptForStorages) String() string { return proto.CompactTextString(m) }
func (*ReceiptForStorages) ProtoMessage()    {}
func (*ReceiptForStorages) Descriptor() ([]byte, []int) {
	return fileDescriptor_0afd3d30283bce23, []int{5}
}

func (m *ReceiptForStorages) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReceiptForStorages.Unmarshal(m, b)
}
func (m *ReceiptForStorages) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReceiptForStorages.Marshal(b, m, deterministic)
}
func (m *ReceiptForStorages) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReceiptForStorages.Merge(m, src)
}
func (m *ReceiptForStorages) XXX_Size() int {
	return xxx_messageInfo_ReceiptForStorages.Size(m)
}
func (m *ReceiptForStorages) XXX_DiscardUnknown() {
	xxx_messageInfo_ReceiptForStorages.DiscardUnknown(m)
}

var xxx_messageInfo_ReceiptForStorages proto.InternalMessageInfo

func (m *ReceiptForStorages) GetReceipts() []*ReceiptForStorage {
	if m != nil {
		return m.Receipts
	}
	return nil
}

type LogForStorage struct {
	Address              []byte   `protobuf:"bytes,1,req,name=address" json:"address,omitempty"`
	Topics               [][]byte `protobuf:"bytes,2,rep,name=topics" json:"topics,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,req,name=data" json:"data,omitempty"`
	BlockNumber          *uint64  `protobuf:"varint,4,req,name=blockNumber" json:"blockNumber,omitempty"`
	TxHash               []byte   `protobuf:"bytes,5,req,name=txHash" json:"txHash,omitempty"`
	TxIndex              *uint32  `protobuf:"varint,6,req,name=txIndex" json:"txIndex,omitempty"`
	BlockHash            []byte   `protobuf:"bytes,7,req,name=blockHash" json:"blockHash,omitempty"`
	Index                *uint32  `protobuf:"varint,8,req,name=index" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogForStorage) Reset()         { *m = LogForStorage{} }
func (m *LogForStorage) String() string { return proto.CompactTextString(m) }
func (*LogForStorage) ProtoMessage()    {}
func (*LogForStorage) Descriptor() ([]byte, []int) {
	return fileDescriptor_0afd3d30283bce23, []int{6}
}

func (m *LogForStorage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogForStorage.Unmarshal(m, b)
}
func (m *LogForStorage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogForStorage.Marshal(b, m, deterministic)
}
func (m *LogForStorage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogForStorage.Merge(m, src)
}
func (m *LogForStorage) XXX_Size() int {
	return xxx_messageInfo_LogForStorage.Size(m)
}
func (m *LogForStorage) XXX_DiscardUnknown() {
	xxx_messageInfo_LogForStorage.DiscardUnknown(m)
}

var xxx_messageInfo_LogForStorage proto.InternalMessageInfo

func (m *LogForStorage) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *LogForStorage) GetTopics() [][]byte {
	if m != nil {
		return m.Topics
	}
	return nil
}

func (m *LogForStorage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *LogForStorage) GetBlockNumber() uint64 {
	if m != nil && m.BlockNumber != nil {
		return *m.BlockNumber
	}
	return 0
}

func (m *LogForStorage) GetTxHash() []byte {
	if m != nil {
		return m.TxHash
	}
	return nil
}

func (m *LogForStorage) GetTxIndex() uint32 {
	if m != nil && m.TxIndex != nil {
		return *m.TxIndex
	}
	return 0
}

func (m *LogForStorage) GetBlockHash() []byte {
	if m != nil {
		return m.BlockHash
	}
	return nil
}

func (m *LogForStorage) GetIndex() uint32 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func init() {
	proto.RegisterType((*BlockHeaderData)(nil), "protobuf.BlockHeaderData")
	proto.RegisterType((*TxData)(nil), "protobuf.TxData")
	proto.RegisterType((*Receipt)(nil), "protobuf.Receipt")
	proto.RegisterType((*Log)(nil), "protobuf.Log")
	proto.RegisterType((*ReceiptForStorage)(nil), "protobuf.ReceiptForStorage")
	proto.RegisterType((*ReceiptForStorages)(nil), "protobuf.ReceiptForStorages")
	proto.RegisterType((*LogForStorage)(nil), "protobuf.LogForStorage")
}

func init() { proto.RegisterFile("protobuf/contract.proto", fileDescriptor_0afd3d30283bce23) }

var fileDescriptor_0afd3d30283bce23 = []byte{
	// 434 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x53, 0xcd, 0x8e, 0xd3, 0x30,
	0x10, 0x56, 0x93, 0x6c, 0x1b, 0x66, 0x5b, 0xad, 0xd6, 0x42, 0x6c, 0x04, 0x1c, 0x42, 0x4e, 0x91,
	0x40, 0x45, 0xe2, 0xc2, 0x99, 0x1f, 0xc1, 0x22, 0x0a, 0x48, 0xde, 0xe5, 0x01, 0x9c, 0xc4, 0x04,
	0x8b, 0x64, 0x1d, 0xd9, 0x0e, 0x0a, 0x2f, 0xc1, 0x23, 0xf0, 0x66, 0xbc, 0x0b, 0xf2, 0xd8, 0x69,
	0x1b, 0x22, 0xc1, 0x05, 0x89, 0x53, 0xf2, 0xcd, 0x78, 0x66, 0xbe, 0xef, 0xf3, 0x18, 0x2e, 0x3a,
	0x25, 0x8d, 0x2c, 0xfa, 0x4f, 0x8f, 0x4b, 0x79, 0x63, 0x14, 0x2b, 0xcd, 0x16, 0x23, 0x24, 0x1e,
	0x13, 0x99, 0x80, 0xb3, 0xe7, 0x8d, 0x2c, 0xbf, 0x5c, 0x72, 0x56, 0x71, 0xf5, 0x92, 0x19, 0x46,
	0x52, 0x38, 0x55, 0xbc, 0xe4, 0xa2, 0x33, 0x97, 0x4c, 0x7f, 0x4e, 0x16, 0x69, 0x90, 0xaf, 0xe9,
	0x71, 0x88, 0xdc, 0x85, 0xb8, 0x66, 0x7a, 0x27, 0x5a, 0x61, 0x92, 0x20, 0x0d, 0xf2, 0x88, 0xee,
	0x31, 0x49, 0x60, 0x55, 0x33, 0xfd, 0x51, 0xf3, 0x2a, 0x09, 0x31, 0x35, 0xc2, 0xec, 0x1a, 0x96,
	0xd7, 0x03, 0x4e, 0xb8, 0x0d, 0x27, 0x9d, 0x12, 0x25, 0xc7, 0xde, 0x11, 0x75, 0xe0, 0x6f, 0x5d,
	0x3b, 0xf6, 0xad, 0x91, 0xcc, 0x75, 0x5d, 0xd3, 0x11, 0x66, 0x3f, 0x16, 0xb0, 0xa2, 0x8e, 0x1b,
	0x79, 0x04, 0xe7, 0x9d, 0xd4, 0xe6, 0xca, 0x30, 0xc3, 0x3f, 0x28, 0xfb, 0xe9, 0xb5, 0xe7, 0x3f,
	0x4f, 0xd8, 0xd3, 0x65, 0xdf, 0xf6, 0x0d, 0x33, 0xe2, 0x2b, 0x7f, 0xed, 0x39, 0xbb, 0xc1, 0xf3,
	0x84, 0xe5, 0x5c, 0x34, 0x52, 0xb6, 0x7e, 0xbe, 0x03, 0xe4, 0x01, 0x44, 0x8d, 0xac, 0x75, 0x12,
	0xa5, 0x61, 0x7e, 0xfa, 0x64, 0xb3, 0x1d, 0x7d, 0xdd, 0xee, 0x64, 0x4d, 0x31, 0x95, 0xbd, 0x85,
	0x70, 0x27, 0x6b, 0xab, 0x80, 0x55, 0x95, 0xe2, 0x7a, 0x64, 0x34, 0x42, 0x72, 0x07, 0x96, 0x46,
	0x76, 0xa2, 0xd4, 0x49, 0x90, 0x86, 0xf9, 0x9a, 0x7a, 0x44, 0x08, 0x44, 0x15, 0x33, 0xcc, 0x0f,
	0xc4, 0xff, 0xec, 0x7b, 0x00, 0xe7, 0x5e, 0xed, 0x2b, 0xa9, 0xae, 0x8c, 0x54, 0xac, 0xe6, 0xff,
	0x41, 0xf7, 0xc3, 0x89, 0xee, 0x8b, 0x89, 0xee, 0x03, 0x31, 0xe7, 0x00, 0x0a, 0x1c, 0x70, 0x97,
	0x4e, 0xb0, 0x87, 0x47, 0xc7, 0xab, 0xb2, 0x9c, 0xac, 0x0a, 0xc9, 0xe1, 0xec, 0x85, 0xdf, 0xd8,
	0x67, 0xde, 0xb4, 0x15, 0x96, 0xfe, 0x1e, 0xce, 0xde, 0x01, 0x99, 0xf9, 0xa1, 0xc9, 0x53, 0x88,
	0xfd, 0xbe, 0x5a, 0x1f, 0x2c, 0xc5, 0x7b, 0x07, 0x8a, 0xb3, 0xf3, 0x74, 0x7f, 0x38, 0xfb, 0xb9,
	0x80, 0xcd, 0x44, 0xc2, 0xbf, 0xb9, 0x37, 0xfb, 0xa6, 0x0a, 0xfb, 0xcc, 0xde, 0xf7, 0x6d, 0xc1,
	0x55, 0x12, 0xa1, 0xdc, 0xe3, 0xd0, 0x9f, 0x4c, 0x32, 0xc3, 0x9b, 0x9b, 0x8a, 0x0f, 0x68, 0xd2,
	0x86, 0x8e, 0x90, 0xdc, 0x87, 0x5b, 0xd8, 0x00, 0x8b, 0x9c, 0x3d, 0x87, 0x80, 0xbd, 0x37, 0x81,
	0x55, 0x31, 0x56, 0x39, 0xf0, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xab, 0x56, 0x4a, 0xcd, 0x12, 0x04,
	0x00, 0x00,
}
