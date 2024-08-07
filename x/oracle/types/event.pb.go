// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: settlus/oracle/v1alpha1/event.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
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

type EventPrevote struct {
	Feeder    string `protobuf:"bytes,1,opt,name=feeder,proto3" json:"feeder,omitempty"`
	Validator string `protobuf:"bytes,2,opt,name=validator,proto3" json:"validator,omitempty"`
	Hash      string `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *EventPrevote) Reset()         { *m = EventPrevote{} }
func (m *EventPrevote) String() string { return proto.CompactTextString(m) }
func (*EventPrevote) ProtoMessage()    {}
func (*EventPrevote) Descriptor() ([]byte, []int) {
	return fileDescriptor_19e458172dc6362a, []int{0}
}
func (m *EventPrevote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventPrevote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventPrevote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventPrevote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventPrevote.Merge(m, src)
}
func (m *EventPrevote) XXX_Size() int {
	return m.Size()
}
func (m *EventPrevote) XXX_DiscardUnknown() {
	xxx_messageInfo_EventPrevote.DiscardUnknown(m)
}

var xxx_messageInfo_EventPrevote proto.InternalMessageInfo

func (m *EventPrevote) GetFeeder() string {
	if m != nil {
		return m.Feeder
	}
	return ""
}

func (m *EventPrevote) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func (m *EventPrevote) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type EventVote struct {
	Feeder    string      `protobuf:"bytes,1,opt,name=feeder,proto3" json:"feeder,omitempty"`
	Validator string      `protobuf:"bytes,2,opt,name=validator,proto3" json:"validator,omitempty"`
	VoteData  []*VoteData `protobuf:"bytes,3,rep,name=vote_data,json=voteData,proto3" json:"vote_data,omitempty"`
}

func (m *EventVote) Reset()         { *m = EventVote{} }
func (m *EventVote) String() string { return proto.CompactTextString(m) }
func (*EventVote) ProtoMessage()    {}
func (*EventVote) Descriptor() ([]byte, []int) {
	return fileDescriptor_19e458172dc6362a, []int{1}
}
func (m *EventVote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventVote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventVote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventVote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventVote.Merge(m, src)
}
func (m *EventVote) XXX_Size() int {
	return m.Size()
}
func (m *EventVote) XXX_DiscardUnknown() {
	xxx_messageInfo_EventVote.DiscardUnknown(m)
}

var xxx_messageInfo_EventVote proto.InternalMessageInfo

func (m *EventVote) GetFeeder() string {
	if m != nil {
		return m.Feeder
	}
	return ""
}

func (m *EventVote) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func (m *EventVote) GetVoteData() []*VoteData {
	if m != nil {
		return m.VoteData
	}
	return nil
}

type EventFeederDelegationConsent struct {
	Feeder    string `protobuf:"bytes,1,opt,name=feeder,proto3" json:"feeder,omitempty"`
	Validator string `protobuf:"bytes,2,opt,name=validator,proto3" json:"validator,omitempty"`
}

func (m *EventFeederDelegationConsent) Reset()         { *m = EventFeederDelegationConsent{} }
func (m *EventFeederDelegationConsent) String() string { return proto.CompactTextString(m) }
func (*EventFeederDelegationConsent) ProtoMessage()    {}
func (*EventFeederDelegationConsent) Descriptor() ([]byte, []int) {
	return fileDescriptor_19e458172dc6362a, []int{2}
}
func (m *EventFeederDelegationConsent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventFeederDelegationConsent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventFeederDelegationConsent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventFeederDelegationConsent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventFeederDelegationConsent.Merge(m, src)
}
func (m *EventFeederDelegationConsent) XXX_Size() int {
	return m.Size()
}
func (m *EventFeederDelegationConsent) XXX_DiscardUnknown() {
	xxx_messageInfo_EventFeederDelegationConsent.DiscardUnknown(m)
}

var xxx_messageInfo_EventFeederDelegationConsent proto.InternalMessageInfo

func (m *EventFeederDelegationConsent) GetFeeder() string {
	if m != nil {
		return m.Feeder
	}
	return ""
}

func (m *EventFeederDelegationConsent) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

type EventOracleConsensusFailed struct {
	ChainId     string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	BlockHeight int64  `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

func (m *EventOracleConsensusFailed) Reset()         { *m = EventOracleConsensusFailed{} }
func (m *EventOracleConsensusFailed) String() string { return proto.CompactTextString(m) }
func (*EventOracleConsensusFailed) ProtoMessage()    {}
func (*EventOracleConsensusFailed) Descriptor() ([]byte, []int) {
	return fileDescriptor_19e458172dc6362a, []int{3}
}
func (m *EventOracleConsensusFailed) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventOracleConsensusFailed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventOracleConsensusFailed.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventOracleConsensusFailed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventOracleConsensusFailed.Merge(m, src)
}
func (m *EventOracleConsensusFailed) XXX_Size() int {
	return m.Size()
}
func (m *EventOracleConsensusFailed) XXX_DiscardUnknown() {
	xxx_messageInfo_EventOracleConsensusFailed.DiscardUnknown(m)
}

var xxx_messageInfo_EventOracleConsensusFailed proto.InternalMessageInfo

func (m *EventOracleConsensusFailed) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

func (m *EventOracleConsensusFailed) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*EventPrevote)(nil), "settlus.oracle.v1alpha1.EventPrevote")
	proto.RegisterType((*EventVote)(nil), "settlus.oracle.v1alpha1.EventVote")
	proto.RegisterType((*EventFeederDelegationConsent)(nil), "settlus.oracle.v1alpha1.EventFeederDelegationConsent")
	proto.RegisterType((*EventOracleConsensusFailed)(nil), "settlus.oracle.v1alpha1.EventOracleConsensusFailed")
}

func init() {
	proto.RegisterFile("settlus/oracle/v1alpha1/event.proto", fileDescriptor_19e458172dc6362a)
}

var fileDescriptor_19e458172dc6362a = []byte{
	// 335 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0x41, 0x4f, 0xea, 0x40,
	0x10, 0xc7, 0xe9, 0xeb, 0x0b, 0x8f, 0x2e, 0x9c, 0xf6, 0xf0, 0x5e, 0x1f, 0x21, 0x0d, 0x54, 0x13,
	0x39, 0xb5, 0x41, 0xef, 0x26, 0x2a, 0x12, 0x3d, 0x69, 0x88, 0x31, 0x86, 0x0b, 0x59, 0xda, 0x91,
	0xdd, 0xb8, 0x76, 0x49, 0x77, 0x68, 0xf4, 0xe8, 0x37, 0xf0, 0x63, 0x79, 0xe4, 0xe8, 0xd1, 0xc0,
	0x17, 0x31, 0x0c, 0x25, 0x9e, 0xb8, 0x70, 0xdb, 0x99, 0xf9, 0xcf, 0x6f, 0xfe, 0x9b, 0x3f, 0x3b,
	0xb0, 0x80, 0xa8, 0xe7, 0x36, 0x36, 0xb9, 0x48, 0x34, 0xc4, 0x45, 0x4f, 0xe8, 0x99, 0x14, 0xbd,
	0x18, 0x0a, 0xc8, 0x30, 0x9a, 0xe5, 0x06, 0x0d, 0xff, 0x57, 0x8a, 0xa2, 0x8d, 0x28, 0xda, 0x8a,
	0x9a, 0x87, 0xbb, 0xb6, 0x4b, 0x21, 0xad, 0x87, 0x0f, 0xac, 0x71, 0xb9, 0xa6, 0xdd, 0xe6, 0x50,
	0x18, 0x04, 0xfe, 0x97, 0x55, 0x1f, 0x01, 0x52, 0xc8, 0x7d, 0xa7, 0xed, 0x74, 0xbd, 0x61, 0x59,
	0xf1, 0x16, 0xf3, 0x0a, 0xa1, 0x55, 0x2a, 0xd0, 0xe4, 0xfe, 0x2f, 0x1a, 0xfd, 0x34, 0x38, 0x67,
	0xbf, 0xa5, 0xb0, 0xd2, 0x77, 0x69, 0x40, 0xef, 0xf0, 0xcd, 0x61, 0x1e, 0xa1, 0xef, 0xf7, 0xe7,
	0x9e, 0x32, 0x6f, 0xed, 0x6a, 0x9c, 0x0a, 0x14, 0xbe, 0xdb, 0x76, 0xbb, 0xf5, 0xe3, 0x4e, 0xb4,
	0xe3, 0xc3, 0xd1, 0xfa, 0x4e, 0x5f, 0xa0, 0x18, 0xd6, 0x8a, 0xf2, 0x15, 0xde, 0xb1, 0x16, 0x59,
	0x18, 0xd0, 0xb1, 0x3e, 0x68, 0x98, 0x0a, 0x54, 0x26, 0xbb, 0x30, 0x99, 0x85, 0x0c, 0xf7, 0x73,
	0x15, 0x8e, 0x58, 0x93, 0xa8, 0x37, 0x64, 0x60, 0xc3, 0xb2, 0x73, 0x3b, 0x10, 0x4a, 0x43, 0xca,
	0xff, 0xb3, 0x5a, 0x22, 0x85, 0xca, 0xc6, 0x2a, 0x2d, 0xa9, 0x7f, 0xa8, 0xbe, 0x4e, 0x79, 0x87,
	0x35, 0x26, 0xda, 0x24, 0x4f, 0x63, 0x09, 0x6a, 0x2a, 0x91, 0xc8, 0xee, 0xb0, 0x4e, 0xbd, 0x2b,
	0x6a, 0x9d, 0x9f, 0x7d, 0x2c, 0x03, 0x67, 0xb1, 0x0c, 0x9c, 0xaf, 0x65, 0xe0, 0xbc, 0xaf, 0x82,
	0xca, 0x62, 0x15, 0x54, 0x3e, 0x57, 0x41, 0x65, 0x74, 0x34, 0x55, 0x28, 0xe7, 0x93, 0x28, 0x31,
	0xcf, 0xf1, 0x36, 0x5a, 0x02, 0xc7, 0x2f, 0xdb, 0x88, 0xf1, 0x75, 0x06, 0x76, 0x52, 0xa5, 0x64,
	0x4f, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0x10, 0x4b, 0xfa, 0x62, 0x3f, 0x02, 0x00, 0x00,
}

func (m *EventPrevote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventPrevote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventPrevote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Feeder) > 0 {
		i -= len(m.Feeder)
		copy(dAtA[i:], m.Feeder)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Feeder)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventVote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventVote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventVote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.VoteData) > 0 {
		for iNdEx := len(m.VoteData) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.VoteData[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintEvent(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Feeder) > 0 {
		i -= len(m.Feeder)
		copy(dAtA[i:], m.Feeder)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Feeder)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventFeederDelegationConsent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventFeederDelegationConsent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventFeederDelegationConsent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Feeder) > 0 {
		i -= len(m.Feeder)
		copy(dAtA[i:], m.Feeder)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Feeder)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventOracleConsensusFailed) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventOracleConsensusFailed) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventOracleConsensusFailed) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BlockHeight != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x10
	}
	if len(m.ChainId) > 0 {
		i -= len(m.ChainId)
		copy(dAtA[i:], m.ChainId)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.ChainId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventPrevote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Feeder)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	return n
}

func (m *EventVote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Feeder)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if len(m.VoteData) > 0 {
		for _, e := range m.VoteData {
			l = e.Size()
			n += 1 + l + sovEvent(uint64(l))
		}
	}
	return n
}

func (m *EventFeederDelegationConsent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Feeder)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	return n
}

func (m *EventOracleConsensusFailed) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChainId)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovEvent(uint64(m.BlockHeight))
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventPrevote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: EventPrevote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventPrevote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feeder", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Feeder = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
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
func (m *EventVote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: EventVote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventVote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feeder", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Feeder = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoteData", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VoteData = append(m.VoteData, &VoteData{})
			if err := m.VoteData[len(m.VoteData)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
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
func (m *EventFeederDelegationConsent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: EventFeederDelegationConsent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventFeederDelegationConsent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feeder", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Feeder = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
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
func (m *EventOracleConsensusFailed) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: EventOracleConsensusFailed: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventOracleConsensusFailed: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
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
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvent = fmt.Errorf("proto: unexpected end of group")
)
