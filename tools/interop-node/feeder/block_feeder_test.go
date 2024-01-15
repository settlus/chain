package feeder_test

import (
	"testing"

	"github.com/settlus/chain/tools/interop-node/feeder"
	oracletypes "github.com/settlus/chain/x/oracle/types"
)

func Test_BlockDataListToBlockDataString(t *testing.T) {
	type args struct {
		bdList []oracletypes.BlockData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				bdList: []oracletypes.BlockData{
					{
						ChainId:     "1",
						BlockNumber: 123,
						BlockHash:   "0x123",
					},
					{
						ChainId:     "2",
						BlockNumber: 456,
						BlockHash:   "0x456",
					},
				},
			},
			want: "1:123:0x123,2:456:0x456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := feeder.BlockDataListToBlockDataString(tt.args.bdList); got != tt.want {
				t.Errorf("BlockDataListToBlockDataString() = %v, want %v", got, tt.want)
			}
		})
	}
}
