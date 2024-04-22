// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: settlus/settlement/v1alpha1/settlement.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_settlus_chain_types "github.com/settlus/chain/types"
	types1 "github.com/settlus/chain/x/oracle/types"
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
	return fileDescriptor_c8003d389c9288b3, []int{0}
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
	return fileDescriptor_c8003d389c9288b3, []int{1}
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
	RequestId  string       `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Recipients []*Recipient `protobuf:"bytes,2,rep,name=recipients,proto3" json:"recipients,omitempty"`
	Amount     types.Coin   `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
	CreatedAt  uint64       `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Nft        *types1.Nft  `protobuf:"bytes,5,opt,name=nft,proto3" json:"nft,omitempty"`
}

func (m *UTXR) Reset()         { *m = UTXR{} }
func (m *UTXR) String() string { return proto.CompactTextString(m) }
func (*UTXR) ProtoMessage()    {}
func (*UTXR) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8003d389c9288b3, []int{2}
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

func (m *UTXR) GetRecipients() []*Recipient {
	if m != nil {
		return m.Recipients
	}
	return nil
}

func (m *UTXR) GetAmount() types.Coin {
	if m != nil {
		return m.Amount
	}
	return types.Coin{}
}

func (m *UTXR) GetCreatedAt() uint64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func (m *UTXR) GetNft() *types1.Nft {
	if m != nil {
		return m.Nft
	}
	return nil
}

type Recipient struct {
	Address github_com_settlus_chain_types.HexAddressString `protobuf:"bytes,1,opt,name=address,proto3,customtype=github.com/settlus/chain/types.HexAddressString" json:"address"`
	Weight  uint32                                          `protobuf:"varint,2,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (m *Recipient) Reset()         { *m = Recipient{} }
func (m *Recipient) String() string { return proto.CompactTextString(m) }
func (*Recipient) ProtoMessage()    {}
func (*Recipient) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8003d389c9288b3, []int{3}
}
func (m *Recipient) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Recipient) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Recipient.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Recipient) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Recipient.Merge(m, src)
}
func (m *Recipient) XXX_Size() int {
	return m.Size()
}
func (m *Recipient) XXX_DiscardUnknown() {
	xxx_messageInfo_Recipient.DiscardUnknown(m)
}

var xxx_messageInfo_Recipient proto.InternalMessageInfo

func (m *Recipient) GetWeight() uint32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "settlus.settlement.v1alpha1.Params")
	proto.RegisterType((*Tenant)(nil), "settlus.settlement.v1alpha1.Tenant")
	proto.RegisterType((*UTXR)(nil), "settlus.settlement.v1alpha1.UTXR")
	proto.RegisterType((*Recipient)(nil), "settlus.settlement.v1alpha1.Recipient")
}

func init() {
	proto.RegisterFile("settlus/settlement/v1alpha1/settlement.proto", fileDescriptor_c8003d389c9288b3)
}

var fileDescriptor_c8003d389c9288b3 = []byte{
	// 586 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xc1, 0x6e, 0xd3, 0x4a,
	0x14, 0x8d, 0xd3, 0x34, 0xaf, 0x9e, 0xbe, 0x02, 0x1a, 0x0a, 0x32, 0x05, 0x9c, 0x28, 0xa0, 0x2a,
	0x08, 0xb0, 0x95, 0xb2, 0xa8, 0x54, 0xb1, 0x69, 0x80, 0x0a, 0x16, 0xa0, 0x30, 0x14, 0x09, 0xb1,
	0x89, 0x26, 0xe3, 0x5b, 0x67, 0x44, 0x3d, 0x63, 0x3c, 0x93, 0xd2, 0xfe, 0x05, 0x4b, 0x96, 0xfd,
	0x04, 0x7e, 0x02, 0xa9, 0xcb, 0x2e, 0x11, 0x8b, 0x0a, 0x35, 0x1b, 0x7e, 0x80, 0x3d, 0xf2, 0xcc,
	0x38, 0x64, 0x03, 0x62, 0x65, 0xdf, 0xe3, 0x33, 0xf7, 0x9c, 0x7b, 0x7d, 0x06, 0xdd, 0x53, 0xa0,
	0xf5, 0xfe, 0x44, 0xc5, 0xe6, 0x09, 0x19, 0x08, 0x1d, 0x1f, 0xf4, 0xe8, 0x7e, 0x3e, 0xa6, 0xbd,
	0x39, 0x2c, 0xca, 0x0b, 0xa9, 0x25, 0xbe, 0xee, 0xd8, 0xd1, 0xdc, 0x97, 0x8a, 0xbd, 0x16, 0x32,
	0xa9, 0x32, 0xa9, 0xe2, 0x11, 0x55, 0x10, 0x1f, 0xf4, 0x46, 0xa0, 0x69, 0x2f, 0x66, 0x92, 0x0b,
	0x7b, 0x78, 0xed, 0x76, 0x25, 0x25, 0x0b, 0xca, 0xf6, 0xe1, 0xb7, 0x8c, 0xad, 0x1d, 0x6b, 0x35,
	0x95, 0xa9, 0x34, 0xaf, 0x71, 0xf9, 0x66, 0xd1, 0xce, 0x67, 0x0f, 0x35, 0x07, 0xb4, 0xa0, 0x99,
	0xc2, 0x0f, 0x91, 0x9f, 0x52, 0x35, 0xcc, 0x0b, 0xce, 0x20, 0xf0, 0xda, 0x5e, 0x77, 0x79, 0xe3,
	0x5a, 0x64, 0xa5, 0xa3, 0x52, 0x3a, 0x72, 0xd2, 0xd1, 0x23, 0xc9, 0x45, 0xbf, 0x71, 0x72, 0xd6,
	0xaa, 0x91, 0xa5, 0x94, 0xaa, 0x41, 0x79, 0x00, 0x8f, 0xd0, 0x15, 0x2b, 0x37, 0xdc, 0x03, 0x18,
	0xe6, 0x50, 0x30, 0x10, 0x9a, 0xa6, 0x10, 0x2c, 0xb4, 0xbd, 0xae, 0xdf, 0x8f, 0x4a, 0xfa, 0xb7,
	0xb3, 0xd6, 0x7a, 0xca, 0xf5, 0x78, 0x32, 0x8a, 0x98, 0xcc, 0x62, 0x37, 0x96, 0x7d, 0xdc, 0x57,
	0xc9, 0xbb, 0x58, 0x1f, 0xe5, 0xa0, 0xa2, 0xc7, 0xc0, 0xc8, 0x65, 0xdb, 0x6c, 0x07, 0x60, 0x30,
	0x6b, 0xb5, 0xd5, 0xf8, 0x74, 0xdc, 0xaa, 0x75, 0xbe, 0x78, 0xa8, 0xb9, 0x0b, 0x82, 0x0a, 0x8d,
	0x2f, 0xa0, 0x3a, 0x4f, 0x8c, 0xd7, 0x06, 0xa9, 0xf3, 0x04, 0x5f, 0x45, 0x4d, 0x9a, 0x64, 0x5c,
	0xa8, 0xa0, 0xde, 0x5e, 0xe8, 0xfa, 0xc4, 0x55, 0x78, 0x15, 0x2d, 0x26, 0x20, 0x64, 0x66, 0xcd,
	0x10, 0x5b, 0xe0, 0x5b, 0x68, 0x25, 0xa7, 0x47, 0x72, 0xa2, 0x4b, 0xbb, 0x5c, 0x26, 0x41, 0xc3,
	0x34, 0xfa, 0xdf, 0x82, 0x03, 0x83, 0xcd, 0x91, 0x32, 0xd0, 0x63, 0x99, 0x04, 0x8b, 0xa6, 0x85,
	0x23, 0x3d, 0x37, 0x18, 0xbe, 0x83, 0x2e, 0x31, 0x29, 0x74, 0x41, 0x99, 0x1e, 0xd2, 0x24, 0x29,
	0x40, 0xa9, 0xa0, 0x69, 0x78, 0x17, 0x2b, 0x7c, 0xdb, 0xc2, 0x5b, 0x4b, 0xe5, 0x0c, 0x3f, 0x8e,
	0x5b, 0x5e, 0xe7, 0xa7, 0x87, 0x1a, 0xaf, 0x77, 0xdf, 0x10, 0x7c, 0x13, 0xa1, 0x02, 0xde, 0x4f,
	0x40, 0xe9, 0xa1, 0x9b, 0xc6, 0x27, 0xbe, 0x43, 0x9e, 0x25, 0x78, 0xa7, 0xfc, 0xcc, 0x78, 0xce,
	0x41, 0x68, 0x3b, 0xd8, 0xf2, 0xc6, 0x7a, 0xf4, 0x97, 0xc0, 0x44, 0xa4, 0xa2, 0x93, 0xb9, 0x93,
	0x78, 0x13, 0x35, 0x69, 0x26, 0x27, 0x42, 0x9b, 0x2d, 0xfc, 0xc3, 0xcf, 0x75, 0xf4, 0xd2, 0x1f,
	0x2b, 0x80, 0x6a, 0x48, 0x86, 0x54, 0xbb, 0x25, 0xf9, 0x0e, 0xd9, 0xd6, 0x38, 0x42, 0x0b, 0x62,
	0x4f, 0x9b, 0xbd, 0x2c, 0x6f, 0xdc, 0x98, 0x19, 0x73, 0xe1, 0x9b, 0x99, 0x7a, 0xb1, 0xa7, 0x49,
	0x49, 0xec, 0x1c, 0x20, 0x7f, 0x66, 0x10, 0xbf, 0x44, 0xff, 0x55, 0x0b, 0x33, 0x83, 0xf7, 0x37,
	0x5d, 0x50, 0xe2, 0xb9, 0xa0, 0x54, 0xf9, 0x66, 0x63, 0xca, 0x85, 0x0b, 0xc9, 0x53, 0x38, 0x74,
	0x2b, 0x7d, 0xa5, 0x0b, 0x2e, 0x52, 0x52, 0xf5, 0x29, 0x43, 0xf0, 0x01, 0x78, 0x3a, 0xd6, 0x41,
	0xbd, 0xed, 0x75, 0x57, 0x88, 0xab, 0xfa, 0x4f, 0x4e, 0xce, 0x43, 0xef, 0xf4, 0x3c, 0xf4, 0xbe,
	0x9f, 0x87, 0xde, 0xc7, 0x69, 0x58, 0x3b, 0x9d, 0x86, 0xb5, 0xaf, 0xd3, 0xb0, 0xf6, 0xf6, 0xee,
	0x1f, 0xb5, 0x0e, 0xe7, 0xaf, 0xaf, 0x11, 0x1e, 0x35, 0xcd, 0xc5, 0x79, 0xf0, 0x2b, 0x00, 0x00,
	0xff, 0xff, 0xab, 0xb3, 0xf7, 0x73, 0xe1, 0x03, 0x00, 0x00,
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
	if m.Nft != nil {
		{
			size, err := m.Nft.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSettlement(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.CreatedAt != 0 {
		i = encodeVarintSettlement(dAtA, i, uint64(m.CreatedAt))
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
	if len(m.Recipients) > 0 {
		for iNdEx := len(m.Recipients) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Recipients[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSettlement(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.RequestId) > 0 {
		i -= len(m.RequestId)
		copy(dAtA[i:], m.RequestId)
		i = encodeVarintSettlement(dAtA, i, uint64(len(m.RequestId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Recipient) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Recipient) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Recipient) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Weight != 0 {
		i = encodeVarintSettlement(dAtA, i, uint64(m.Weight))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.Address.Size()
		i -= size
		if _, err := m.Address.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSettlement(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
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
	if len(m.Recipients) > 0 {
		for _, e := range m.Recipients {
			l = e.Size()
			n += 1 + l + sovSettlement(uint64(l))
		}
	}
	l = m.Amount.Size()
	n += 1 + l + sovSettlement(uint64(l))
	if m.CreatedAt != 0 {
		n += 1 + sovSettlement(uint64(m.CreatedAt))
	}
	if m.Nft != nil {
		l = m.Nft.Size()
		n += 1 + l + sovSettlement(uint64(l))
	}
	return n
}

func (m *Recipient) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Address.Size()
	n += 1 + l + sovSettlement(uint64(l))
	if m.Weight != 0 {
		n += 1 + sovSettlement(uint64(m.Weight))
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
				return fmt.Errorf("proto: wrong wireType = %d for field Recipients", wireType)
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
			m.Recipients = append(m.Recipients, &Recipient{})
			if err := m.Recipients[len(m.Recipients)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			m.CreatedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatedAt |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nft", wireType)
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
			if m.Nft == nil {
				m.Nft = &types1.Nft{}
			}
			if err := m.Nft.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *Recipient) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Recipient: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Recipient: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			if err := m.Address.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			m.Weight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSettlement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Weight |= uint32(b&0x7F) << shift
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
