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
		dataString []string
		salt       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single chain",
			args: args{
				dataString: []string{"1/0x123/0x1:0xfoo"},
				salt:       "foo",
			},
			want: "A13CCB8173F8C6919BAE51EAED94FA1066BCAF513F5C5286DBA8808A73C63449",
		}, {
			name: "multiple chains",
			args: args{
				dataString: []string{"1/0x123/0x1:0xfoo", "2/0x456/0x2:0xbar"},
				salt:       "bar",
			},
			want: "5669721A5163DA85815FB766730AC8755D2AC5AF053373EC3CB67135C3BCE21C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			voteData := types.VoteDataArr{
				{
					Topic: oracletypes.OracleTopic_OWNERSHIP,
					Data:  tt.args.dataString,
				},
			}
			if got := feeder.GeneratePrevoteHash(voteData, tt.args.salt); got != tt.want {
				t.Errorf("GeneratePrevoteHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
