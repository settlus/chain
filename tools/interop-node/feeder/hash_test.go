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
				blockDataString: []string{"1/0x123/0x1"},
				salt:            "foo",
			},
			want: "8D09B5CA5B1208D79CD5A507187E762EE32B7403ECE64F313B6F9F071D10AEB2",
		}, {
			name: "multiple chains",
			args: args{
				blockDataString: []string{"1/0x123/0x1", "2/0x456/0x2"},
				salt:            "bar",
			},
			want: "7E36F2C54F6354F0B1F2D0B2FE71AF21F131D08E2B092A03840CAC9EEC47BF01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			voteData := types.VoteDataArr{
				{
					Topic: oracletypes.OracleTopic_OWNERSHIP,
					Data:  tt.args.blockDataString,
				},
			}
			if got := feeder.GeneratePrevoteHash(voteData, tt.args.salt); got != tt.want {
				t.Errorf("GeneratePrevoteHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
