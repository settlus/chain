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
				t.Errorf("PadBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateHexString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid even hex string",
			args: args{
				s: "0x01",
			},
			want: true,
		},
		{
			name: "valid odd hex string",
			args: args{
				s: "0x0",
			},
			want: false,
		},
		{
			name: "invalid hex string",
			args: args{
				s: "0x0g",
			},
			want: false,
		},
		{
			name: "invalid hex string",
			args: args{
				s: "x1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			types.ValidateHexString(tt.args.s)
		})
	}
}
