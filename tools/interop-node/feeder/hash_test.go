package feeder_test

import (
	"testing"

	"github.com/settlus/chain/tools/interop-node/feeder"
	"github.com/settlus/chain/tools/interop-node/types"
	oracletypes "github.com/settlus/chain/x/oracle/types"
)

func Test_GenerateSalt(t *testing.T) {
	salts := make(map[string]bool)
	for i := 0; i < 100; i++ {
		salt, err := feeder.GenerateSalt()
		if err != nil {
			t.Errorf("GenerateSalt() error = %v", err)
		}
		if _, ok := salts[salt]; ok {
			t.Errorf("salt %s already exists", salt)
		}
		salts[salt] = true
	}
}

func Test_GeneratePrevoteHash(t *testing.T) {
	type args struct {
		blockDataString []string
		salt            string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single chain",
			args: args{
				blockDataString: []string{"1:123:0x123"},
				salt:            "foo",
			},
			want: "01A164031A468DE61267F23A6BD7642AA33422C983D0E298085AEE1244A51F40",
		}, {
			name: "multiple chains",
			args: args{
				blockDataString: []string{"1:123:0x123", "2:456:0x456"},
				salt:            "bar",
			},
			want: "1776F5F1BCACEFC9E75DA6623C9C9B1AA6DDF9831DBDEC3453D0C69B380FBE97",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			voteData := types.VoteDataArr{
				{
					Topic: oracletypes.OracleTopic_BLOCK,
					Data:  tt.args.blockDataString,
				},
			}
			if got := feeder.GeneratePrevoteHash(voteData, tt.args.salt); got != tt.want {
				t.Errorf("GeneratePrevoteHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
