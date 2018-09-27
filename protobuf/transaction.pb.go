// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/transaction.proto

package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Transactions struct {
	Txs                  []*Transaction `protobuf:"bytes,1,rep,name=txs" json:"txs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Transactions) Reset()         { *m = Transactions{} }
func (m *Transactions) String() string { return proto.CompactTextString(m) }
func (*Transactions) ProtoMessage()    {}
func (*Transactions) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f5cb32b462e88af0, []int{0}
}
func (m *Transactions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transactions.Unmarshal(m, b)
}
func (m *Transactions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transactions.Marshal(b, m, deterministic)
}
func (dst *Transactions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transactions.Merge(dst, src)
}
func (m *Transactions) XXX_Size() int {
	return xxx_messageInfo_Transactions.Size(m)
}
func (m *Transactions) XXX_DiscardUnknown() {
	xxx_messageInfo_Transactions.DiscardUnknown(m)
}

var xxx_messageInfo_Transactions proto.InternalMessageInfo

func (m *Transactions) GetTxs() []*Transaction {
	if m != nil {
		return m.Txs
	}
	return nil
}

type Transaction struct {
	Version              *uint32          `protobuf:"varint,1,req,name=version" json:"version,omitempty"`
	From                 *TransactionPeer `protobuf:"bytes,2,req,name=from" json:"from,omitempty"`
	To                   *TransactionPeer `protobuf:"bytes,3,req,name=to" json:"to,omitempty"`
	Amount               *Amount          `protobuf:"bytes,4,req,name=amount" json:"amount,omitempty"`
	Time                 *int64           `protobuf:"varint,5,req,name=time" json:"time,omitempty"`
	Nounce               *uint32          `protobuf:"varint,6,req,name=nounce" json:"nounce,omitempty"`
	Extra                []byte           `protobuf:"bytes,7,opt,name=extra" json:"extra,omitempty"`
	Sign                 []byte           `protobuf:"bytes,8,opt,name=sign" json:"sign,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f5cb32b462e88af0, []int{1}
}
func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (dst *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(dst, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetVersion() uint32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func (m *Transaction) GetFrom() *TransactionPeer {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Transaction) GetTo() *TransactionPeer {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Transaction) GetAmount() *Amount {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *Transaction) GetTime() int64 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

func (m *Transaction) GetNounce() uint32 {
	if m != nil && m.Nounce != nil {
		return *m.Nounce
	}
	return 0
}

func (m *Transaction) GetExtra() []byte {
	if m != nil {
		return m.Extra
	}
	return nil
}

func (m *Transaction) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

type TransactionPeer struct {
	AccountID            *AccountID `protobuf:"bytes,1,req,name=accountID" json:"accountID,omitempty"`
	Extra                []byte     `protobuf:"bytes,2,opt,name=extra" json:"extra,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *TransactionPeer) Reset()         { *m = TransactionPeer{} }
func (m *TransactionPeer) String() string { return proto.CompactTextString(m) }
func (*TransactionPeer) ProtoMessage()    {}
func (*TransactionPeer) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f5cb32b462e88af0, []int{2}
}
func (m *TransactionPeer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionPeer.Unmarshal(m, b)
}
func (m *TransactionPeer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionPeer.Marshal(b, m, deterministic)
}
func (dst *TransactionPeer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionPeer.Merge(dst, src)
}
func (m *TransactionPeer) XXX_Size() int {
	return xxx_messageInfo_TransactionPeer.Size(m)
}
func (m *TransactionPeer) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionPeer.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionPeer proto.InternalMessageInfo

func (m *TransactionPeer) GetAccountID() *AccountID {
	if m != nil {
		return m.AccountID
	}
	return nil
}

func (m *TransactionPeer) GetExtra() []byte {
	if m != nil {
		return m.Extra
	}
	return nil
}

type AccountID struct {
	Id                   []byte   `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountID) Reset()         { *m = AccountID{} }
func (m *AccountID) String() string { return proto.CompactTextString(m) }
func (*AccountID) ProtoMessage()    {}
func (*AccountID) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f5cb32b462e88af0, []int{3}
}
func (m *AccountID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountID.Unmarshal(m, b)
}
func (m *AccountID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountID.Marshal(b, m, deterministic)
}
func (dst *AccountID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountID.Merge(dst, src)
}
func (m *AccountID) XXX_Size() int {
	return xxx_messageInfo_AccountID.Size(m)
}
func (m *AccountID) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountID.DiscardUnknown(m)
}

var xxx_messageInfo_AccountID proto.InternalMessageInfo

func (m *AccountID) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

type Hash struct {
	Data                 []byte   `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hash) Reset()         { *m = Hash{} }
func (m *Hash) String() string { return proto.CompactTextString(m) }
func (*Hash) ProtoMessage()    {}
func (*Hash) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f5cb32b462e88af0, []int{4}
}
func (m *Hash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hash.Unmarshal(m, b)
}
func (m *Hash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hash.Marshal(b, m, deterministic)
}
func (dst *Hash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hash.Merge(dst, src)
}
func (m *Hash) XXX_Size() int {
	return xxx_messageInfo_Hash.Size(m)
}
func (m *Hash) XXX_DiscardUnknown() {
	xxx_messageInfo_Hash.DiscardUnknown(m)
}

var xxx_messageInfo_Hash proto.InternalMessageInfo

func (m *Hash) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Amount struct {
	Value                *int32   `protobuf:"varint,1,req,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Amount) Reset()         { *m = Amount{} }
func (m *Amount) String() string { return proto.CompactTextString(m) }
func (*Amount) ProtoMessage()    {}
func (*Amount) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_f5cb32b462e88af0, []int{5}
}
func (m *Amount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Amount.Unmarshal(m, b)
}
func (m *Amount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Amount.Marshal(b, m, deterministic)
}
func (dst *Amount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Amount.Merge(dst, src)
}
func (m *Amount) XXX_Size() int {
	return xxx_messageInfo_Amount.Size(m)
}
func (m *Amount) XXX_DiscardUnknown() {
	xxx_messageInfo_Amount.DiscardUnknown(m)
}

var xxx_messageInfo_Amount proto.InternalMessageInfo

func (m *Amount) GetValue() int32 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*Transactions)(nil), "protobuf.Transactions")
	proto.RegisterType((*Transaction)(nil), "protobuf.Transaction")
	proto.RegisterType((*TransactionPeer)(nil), "protobuf.TransactionPeer")
	proto.RegisterType((*AccountID)(nil), "protobuf.AccountID")
	proto.RegisterType((*Hash)(nil), "protobuf.Hash")
	proto.RegisterType((*Amount)(nil), "protobuf.Amount")
}

func init() {
	proto.RegisterFile("protobuf/transaction.proto", fileDescriptor_transaction_f5cb32b462e88af0)
}

var fileDescriptor_transaction_f5cb32b462e88af0 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0xbf, 0x4f, 0xfb, 0x30,
	0x10, 0xc5, 0x55, 0x27, 0x4d, 0xdb, 0x4b, 0xbf, 0x5f, 0xd0, 0xf1, 0x43, 0xa6, 0x48, 0x28, 0xca,
	0x42, 0x18, 0x28, 0xa2, 0x0b, 0x33, 0x12, 0x03, 0x6c, 0xc8, 0x62, 0x62, 0x33, 0xa9, 0x0b, 0x96,
	0xa8, 0x8d, 0x6c, 0xa7, 0xea, 0xca, 0x7f, 0x8e, 0x72, 0x4d, 0x48, 0x84, 0x90, 0xd8, 0xee, 0xde,
	0xfb, 0xe8, 0xdd, 0xb3, 0x61, 0xf6, 0xe1, 0x6c, 0xb0, 0x2f, 0xd5, 0xea, 0x2a, 0x38, 0x69, 0xbc,
	0x2c, 0x83, 0xb6, 0x66, 0x4e, 0x22, 0x8e, 0x5b, 0x2f, 0xbf, 0x81, 0xe9, 0x53, 0x67, 0x7b, 0x3c,
	0x87, 0x28, 0x6c, 0x3d, 0x1f, 0x64, 0x51, 0x91, 0x2e, 0x8e, 0xe6, 0x2d, 0x37, 0xef, 0x41, 0xa2,
	0x26, 0xf2, 0x4f, 0x06, 0x69, 0x4f, 0x44, 0x0e, 0xa3, 0x8d, 0x72, 0x5e, 0x5b, 0xc3, 0x07, 0x19,
	0x2b, 0xfe, 0x89, 0x76, 0xc5, 0x4b, 0x88, 0x57, 0xce, 0xae, 0x39, 0xcb, 0x58, 0x91, 0x2e, 0x4e,
	0x7e, 0xcd, 0x7c, 0x54, 0xca, 0x09, 0xc2, 0xf0, 0x02, 0x58, 0xb0, 0x3c, 0xfa, 0x0b, 0x66, 0xc1,
	0x62, 0x01, 0x89, 0x5c, 0xdb, 0xca, 0x04, 0x1e, 0x13, 0xbe, 0xdf, 0xe1, 0xb7, 0xa4, 0x8b, 0xc6,
	0x47, 0x84, 0x38, 0xe8, 0xb5, 0xe2, 0xc3, 0x8c, 0x15, 0x91, 0xa0, 0x19, 0x8f, 0x21, 0x31, 0xb6,
	0x32, 0xa5, 0xe2, 0x09, 0x15, 0x6e, 0x36, 0x3c, 0x84, 0xa1, 0xda, 0x06, 0x27, 0xf9, 0x28, 0x1b,
	0x14, 0x53, 0xb1, 0x5b, 0xea, 0x04, 0xaf, 0x5f, 0x0d, 0x1f, 0x93, 0x48, 0x73, 0xfe, 0x0c, 0x7b,
	0x3f, 0x6a, 0xe1, 0x35, 0x4c, 0x64, 0x59, 0xd6, 0x37, 0x1f, 0xee, 0xe8, 0x23, 0xd2, 0xc5, 0x41,
	0xaf, 0x55, 0x6b, 0x89, 0x8e, 0xea, 0xee, 0xb1, 0xde, 0xbd, 0xfc, 0x14, 0x26, 0xdf, 0x34, 0xfe,
	0x07, 0xa6, 0x97, 0x14, 0x37, 0x15, 0x4c, 0x2f, 0xf3, 0x19, 0xc4, 0xf7, 0xd2, 0xbf, 0xd5, 0xa5,
	0x96, 0x32, 0xc8, 0xc6, 0xa1, 0x39, 0x3f, 0x83, 0x64, 0xf7, 0xf8, 0x3a, 0x78, 0x23, 0xdf, 0x2b,
	0x45, 0xf6, 0x50, 0xec, 0x96, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2f, 0xa2, 0x9a, 0x1a, 0x18,
	0x02, 0x00, 0x00,
}