// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: settlus/settlement/settlement.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_settlus_chain_types "github.com/settlus/chain/types"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Params defines the parameters for the module.
type Params struct {
	GasPrice            types.Coin                             `protobuf:"bytes,1,opt,name=gas_price,json=gasPrice,proto3" json:"gas_price"`
	OracleFeePercentage github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=oracle_fee_percentage,json=oracleFeePercentage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"oracle_fee_percentage"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_e83090d94702b861, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetGasPrice() types.Coin {
	if m != nil {
		return m.GasPrice
	}
	return types.Coin{}
}

// Tenant defines the tenant parameters.
type Tenant struct {
	Id              uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Admins          []string `protobuf:"bytes,2,rep,name=admins,proto3" json:"admins,omitempty"`
	Denom           string   `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
	PayoutPeriod    uint64   `protobuf:"varint,4,opt,name=payout_period,json=payoutPeriod,proto3" json:"payout_period,omitempty"`
	PayoutMethod    string   `protobuf:"bytes,5,opt,name=payout_method,json=payoutMethod,proto3" json:"payout_method,omitempty"`
	ContractAddress string   `protobuf:"bytes,6,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
}

func (m *Tenant) Reset()      { *m = Tenant{} }
func (*Tenant) ProtoMessage() {}
func (*Tenant) Descriptor() ([]byte, []int) {
	return fileDescriptor_e83090d94702b861, []int{1}
}
func (m *Tenant) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Tenant) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Tenant.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Tenant) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tenant.Merge(m, src)
}
func (m *Tenant) XXX_Size() int {
	return m.Size()
}
func (m *Tenant) XXX_DiscardUnknown() {
	xxx_messageInfo_Tenant.DiscardUnknown(m)
}

var xxx_messageInfo_Tenant proto.InternalMessageInfo

func (m *Tenant) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Tenant) GetAdmins() []string {
	if m != nil {
		return m.Admins
	}
	return nil
}

func (m *Tenant) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *Tenant) GetPayoutPeriod() uint64 {
	if m != nil {
		return m.PayoutPeriod
	}
	return 0
}

func (m *Tenant) GetPayoutMethod() string {
	if m != nil {
		return m.PayoutMethod
	}
	return ""
}

func (m *Tenant) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

// UTXR defines the unspent transaction record.
type UTXR struct {
	RequestId   string                                          `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Recipient   github_com_settlus_chain_types.HexAddressString `protobuf:"bytes,2,opt,name=recipient,proto3,customtype=github.com/settlus/chain/types.HexAddressString" json:"recipient"`
	Amount      types.Coin                                      `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
	PayoutBlock uint64                                          `protobuf:"varint,4,opt,name=payout_block,json=payoutBlock,proto3" json:"payout_block,omitempty"`
}

func (m *UTXR) Reset()         { *m = UTXR{} }
func (m *UTXR) String() string { return proto.CompactTextString(m) }
func (*UTXR) ProtoMessage()    {}
func (*UTXR) Descriptor() ([]byte, []int) {
	return fileDescriptor_e83090d94702b861, []int{2}
}
func (m *UTXR) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UTXR) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UTXR.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UTXR) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UTXR.Merge(m, src)
}
func (m *UTXR) XXX_Size() int {
	return m.Size()
}
func (m *UTXR) XXX_DiscardUnknown() {
	xxx_messageInfo_UTXR.DiscardUnknown(m)
}

var xxx_messageInfo_UTXR proto.InternalMessageInfo

func (m *UTXR) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

func (m *UTXR) GetAmount() types.Coin {
	if m != nil {
		return m.Amount
	}
	return types.Coin{}
}

func (m *UTXR) GetPayoutBlock() uint64 {
	if m != nil {
		return m.PayoutBlock
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "settlus.settlement.Params")
	proto.RegisterType((*Tenant)(nil), "settlus.settlement.Tenant")
	proto.RegisterType((*UTXR)(nil), "settlus.settlement.UTXR")
}

func init() {
	proto.RegisterFile("settlus/settlement/settlement.proto", fileDescriptor_e83090d94702b861)
}

var fileDescriptor_e83090d94702b861 = []byte{
	// 504 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x4d, 0x6b, 0xd4, 0x40,
	0x18, 0xde, 0x6c, 0xb7, 0xa1, 0x99, 0xfa, 0xc5, 0x58, 0x25, 0x16, 0xcc, 0xae, 0x5b, 0x90, 0x15,
	0x31, 0xa1, 0x7a, 0x28, 0x14, 0x2f, 0xae, 0x1f, 0xe8, 0x41, 0x58, 0x62, 0x0b, 0xe2, 0x25, 0x4c,
	0x26, 0xaf, 0xd9, 0xa1, 0x9b, 0x99, 0x38, 0x33, 0x91, 0xf6, 0x5f, 0x78, 0xf4, 0xd8, 0x9f, 0xe0,
	0x9f, 0x10, 0x7a, 0xec, 0xb1, 0x78, 0x28, 0xb2, 0x7b, 0xf1, 0x67, 0x48, 0x66, 0xa6, 0x76, 0x2f,
	0x42, 0x4f, 0xf3, 0xce, 0x93, 0xe7, 0x7d, 0xf2, 0xbc, 0xcf, 0xbc, 0x68, 0x4b, 0x81, 0xd6, 0xb3,
	0x46, 0x25, 0xe6, 0x84, 0x0a, 0xb8, 0x5e, 0x2a, 0xe3, 0x5a, 0x0a, 0x2d, 0x30, 0x76, 0xa4, 0xf8,
	0xf2, 0xcb, 0x66, 0x44, 0x85, 0xaa, 0x84, 0x4a, 0x72, 0xa2, 0x20, 0xf9, 0xba, 0x9d, 0x83, 0x26,
	0xdb, 0x09, 0x15, 0x8c, 0xdb, 0x9e, 0xcd, 0x8d, 0x52, 0x94, 0xc2, 0x94, 0x49, 0x5b, 0x59, 0x74,
	0xf8, 0xc3, 0x43, 0xfe, 0x84, 0x48, 0x52, 0x29, 0xfc, 0x1c, 0x05, 0x25, 0x51, 0x59, 0x2d, 0x19,
	0x85, 0xd0, 0x1b, 0x78, 0xa3, 0xf5, 0xa7, 0xf7, 0x62, 0x2b, 0x1a, 0xb7, 0xa2, 0xb1, 0x13, 0x8d,
	0x5f, 0x0a, 0xc6, 0xc7, 0xbd, 0x93, 0xf3, 0x7e, 0x27, 0x5d, 0x2b, 0x89, 0x9a, 0xb4, 0x0d, 0x38,
	0x47, 0x77, 0x84, 0x24, 0x74, 0x06, 0xd9, 0x67, 0x80, 0xac, 0x06, 0x49, 0x81, 0x6b, 0x52, 0x42,
	0xb8, 0x32, 0xf0, 0x46, 0xc1, 0x38, 0x6e, 0xe9, 0xbf, 0xce, 0xfb, 0x0f, 0x4b, 0xa6, 0xa7, 0x4d,
	0x1e, 0x53, 0x51, 0x25, 0xce, 0xb0, 0x3d, 0x9e, 0xa8, 0xe2, 0x20, 0xd1, 0x47, 0x35, 0xa8, 0xf8,
	0x15, 0xd0, 0xf4, 0xb6, 0x15, 0x7b, 0x03, 0x30, 0xf9, 0x27, 0xb5, 0xdb, 0xfb, 0x7e, 0xdc, 0xef,
	0x0c, 0x7f, 0x7a, 0xc8, 0xdf, 0x03, 0x4e, 0xb8, 0xc6, 0x37, 0x50, 0x97, 0x15, 0xc6, 0x6b, 0x2f,
	0xed, 0xb2, 0x02, 0xdf, 0x45, 0x3e, 0x29, 0x2a, 0xc6, 0x55, 0xd8, 0x1d, 0xac, 0x8c, 0x82, 0xd4,
	0xdd, 0xf0, 0x06, 0x5a, 0x2d, 0x80, 0x8b, 0xca, 0x9a, 0x49, 0xed, 0x05, 0x6f, 0xa1, 0xeb, 0x35,
	0x39, 0x12, 0x8d, 0x6e, 0xed, 0x32, 0x51, 0x84, 0x3d, 0x23, 0x74, 0xcd, 0x82, 0x13, 0x83, 0x2d,
	0x91, 0x2a, 0xd0, 0x53, 0x51, 0x84, 0xab, 0x46, 0xc2, 0x91, 0xde, 0x1b, 0x0c, 0x3f, 0x42, 0xb7,
	0xa8, 0xe0, 0x5a, 0x12, 0xaa, 0x33, 0x52, 0x14, 0x12, 0x94, 0x0a, 0x7d, 0xc3, 0xbb, 0x79, 0x81,
	0xbf, 0xb0, 0xf0, 0xee, 0x5a, 0x3b, 0xc3, 0x9f, 0xe3, 0xbe, 0x37, 0x3c, 0xf3, 0x50, 0x6f, 0x7f,
	0xef, 0x63, 0x8a, 0xef, 0x23, 0x24, 0xe1, 0x4b, 0x03, 0x4a, 0x67, 0x6e, 0x9a, 0x20, 0x0d, 0x1c,
	0xf2, 0xae, 0xc0, 0xfb, 0x28, 0x90, 0x40, 0x59, 0xcd, 0x80, 0xeb, 0xb0, 0x6b, 0xd2, 0xdc, 0x71,
	0x69, 0x26, 0x4b, 0x69, 0x5e, 0xec, 0x0d, 0x9d, 0x12, 0xc6, 0x5d, 0x92, 0x6f, 0xe1, 0xd0, 0xfd,
	0xf7, 0x83, 0x96, 0x8c, 0x97, 0xe9, 0xa5, 0x12, 0xde, 0x41, 0x3e, 0xa9, 0x44, 0xc3, 0xb5, 0x09,
	0xe5, 0x0a, 0x6f, 0xed, 0xe8, 0xf8, 0x01, 0x72, 0xc3, 0x67, 0xf9, 0x4c, 0xd0, 0x03, 0x97, 0xda,
	0xba, 0xc5, 0xc6, 0x2d, 0x34, 0x7e, 0x7d, 0x32, 0x8f, 0xbc, 0xd3, 0x79, 0xe4, 0xfd, 0x9e, 0x47,
	0xde, 0xb7, 0x45, 0xd4, 0x39, 0x5d, 0x44, 0x9d, 0xb3, 0x45, 0xd4, 0xf9, 0xf4, 0xf8, 0xbf, 0x8e,
	0x0f, 0x97, 0x37, 0xde, 0xd8, 0xcf, 0x7d, 0xb3, 0xa3, 0xcf, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff,
	0xff, 0x88, 0x35, 0xa2, 0x14, 0x03, 0x00, 0x00,
}

func (this *Tenant) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Tenant)
	if !ok {
		that2, ok := that.(Tenant)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Id != that1.Id {
		return false
	}
	if len(this.Admins) != len(that1.Admins) {
		return false
	}
	for i := range this.Admins {
		if this.Admins[i] != that1.Admins[i] {
			return false
		}
	}
	if this.Denom != that1.Denom {
		return false
	}
	if this.PayoutPeriod != that1.PayoutPeriod {
		return false
	}
	if this.PayoutMethod != that1.PayoutMethod {
		return false
	}
	if this.ContractAddress != that1.ContractAddress {
		return false
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.OracleFeePercentage.Size()
		i -= size
		if _, err := m.OracleFeePercentage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSettlement(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.GasPrice.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintSettlement(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Tenant) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Tenant) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Tenant) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ContractAddress) > 0 {
		i -= len(m.ContractAddress)
		copy(dAtA[i:], m.ContractAddress)
		i = encodeVarintSettlement(dAtA, i, uint64(len(m.ContractAddress)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.PayoutMethod) > 0 {
		i -= len(m.PayoutMethod)
		copy(dAtA[i:], m.PayoutMethod)
		i = encodeVarintSettlement(dAtA, i, uint64(len(m.PayoutMethod)))
		i--
		dAtA[i] = 0x2a
	}
	if m.PayoutPeriod != 0 {
		i = encodeVarintSettlement(dAtA, i, uint64(m.PayoutPeriod))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintSettlement(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Admins) > 0 {
		for iNdEx := len(m.Admins) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Admins[iNdEx])
			copy(dAtA[i:], m.Admins[iNdEx])
			i = encodeVarintSettlement(dAtA, i, uint64(len(m.Admins[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Id != 0 {
		i = encodeVarintSettlement(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *UTXR) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UTXR) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UTXR) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PayoutBlock != 0 {
		i = encodeVarintSettlement(dAtA, i, uint64(m.PayoutBlock))
		i--
		dAtA[i] = 0x20
	}
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintSettlement(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Recipient.Size()
		i -= size
		if _, err := m.Recipient.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSettlement(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.RequestId) > 0 {
		i -= len(m.RequestId)
		copy(dAtA[i:], m.RequestId)
		i = encodeVarintSettlement(dAtA, i, uint64(len(m.RequestId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSettlement(dAtA []byte, offset int, v uint64) int {
	offset -= sovSettlement(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.GasPrice.Size()
	n += 1 + l + sovSettlement(uint64(l))
	l = m.OracleFeePercentage.Size()
	n += 1 + l + sovSettlement(uint64(l))
	return n
}

func (m *Tenant) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovSettlement(uint64(m.Id))
	}
	if len(m.Admins) > 0 {
		for _, s := range m.Admins {
			l = len(s)
			n += 1 + l + sovSettlement(uint64(l))
		}
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovSettlement(uint64(l))
	}
	if m.PayoutPeriod != 0 {
		n += 1 + sovSettlement(uint64(m.PayoutPeriod))
	}
	l = len(m.PayoutMethod)
	if l > 0 {
		n += 1 + l + sovSettlement(uint64(l))
	}
	l = len(m.ContractAddress)
	if l > 0 {
		n += 1 + l + sovSettlement(uint64(l))
	}
	return n
}

func (m *UTXR) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RequestId)
	if l > 0 {
		n += 1 + l + sovSettlement(uint64(l))
	}
	l = m.Recipient.Size()
	n += 1 + l + sovSettlement(uint64(l))
	l = m.Amount.Size()
	n += 1 + l + sovSettlement(uint64(l))
	if m.PayoutBlock != 0 {
		n += 1 + sovSettlement(uint64(m.PayoutBlock))
	}
	return n
}

func sovSettlement(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSettlement(x uint64) (n int) {
	return sovSettlement(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSettlement
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasPrice", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSettlement
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSettlement
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GasPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OracleFeePercentage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSettlement
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSettlement
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OracleFeePercentage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSettlement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSettlement
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Tenant) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSettlement
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Tenant: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Tenant: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Admins", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSettlement
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSettlement
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Admins = append(m.Admins, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSettlement
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSettlement
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PayoutPeriod", wireType)
			}
			m.PayoutPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PayoutPeriod |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PayoutMethod", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSettlement
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSettlement
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PayoutMethod = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSettlement
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSettlement
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSettlement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSettlement
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UTXR) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSettlement
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UTXR: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UTXR: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSettlement
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSettlement
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RequestId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSettlement
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSettlement
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Recipient.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSettlement
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSettlement
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PayoutBlock", wireType)
			}
			m.PayoutBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PayoutBlock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSettlement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSettlement
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSettlement(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSettlement
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSettlement
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSettlement
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSettlement
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSettlement        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSettlement          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSettlement = fmt.Errorf("proto: unexpected end of group")
)
