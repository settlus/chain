package signer

import (
	"context"
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	configtypes "github.com/settlus/chain/tools/interop-node/config"
)

type Signer interface {
	PubKey() cryptotypes.PubKey
	Sign(data []byte) ([]byte, error)
}

func NewSigner(ctx context.Context, config *configtypes.Config) Signer {
	switch config.Feeder.SignerMode {
	case configtypes.AwsKms:
		return NewKmsSigner(ctx, config.Feeder.Key)
	case configtypes.Local:
		return NewLocalSigner(config.Feeder.Key)
	default:
		panic(fmt.Sprintf("invalid signer mode, must be one of: %s, %s", configtypes.AwsKms, configtypes.Local))
	}
}
