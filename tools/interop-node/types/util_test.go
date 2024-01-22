package types_test

import (
	"bytes"
	"testing"

	"github.com/settlus/chain/tools/interop-node/types"
)

func Test_PadBytes(t *testing.T) {
	type args struct {
		pad int
		b   []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "pad 32 bytes",
			args: args{
				pad: 32,
				b:   []byte{0x12, 0x01},
			},
			want: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x01,
			},
		},
		{
			name: "pad 20 bytes",
			args: args{
				pad: 20,
				b:   []byte{0x01},
			},
			want: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x01,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := types.PadBytes(tt.args.pad, tt.args.b); !bytes.Equal(got, tt.want) {
				t.Errorf("PadBytes() = %v, ok %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateHexString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		ok   bool
	}{
		{
			name: "valid even hex string",
			s:    "0x01",
			ok:   true,
		},
		{
			name: "valid odd hex string",
			s:    "0x0",
			ok:   false,
		},
		{
			name: "invalid hex string",
			s:    "0x0g",
			ok:   false,
		},
		{
			name: "invalid hex string",
			s:    "x1",
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := types.ValidateHexString(tt.s); ok != tt.ok {
				t.Errorf("ValidateHexString() = %v, ok %v", types.ValidateHexString(tt.s), tt.ok)
			}
		})
	}
}
