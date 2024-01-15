package feeder_test

import (
	"testing"

	"github.com/settlus/chain/tools/interop-node/feeder"
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
		blockDataString string
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
				blockDataString: "1:123:0x123",
				salt:            "foo",
			},
			want: "01A164031A468DE61267F23A6BD7642AA33422C983D0E298085AEE1244A51F40",
		}, {
			name: "multiple chains",
			args: args{
				blockDataString: "1:123:0x123,2:456:0x456",
				salt:            "bar",
			},
			want: "701589DAB7F4883BFD8CCC3AE2A4A90D30C5AE3C21486BB0C346D47A02B4FE05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := feeder.GeneratePrevoteHash(tt.args.blockDataString, tt.args.salt); got != tt.want {
				t.Errorf("GeneratePrevoteHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
