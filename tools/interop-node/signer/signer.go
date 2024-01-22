package signer

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	configtypes "github.com/settlus/chain/tools/interop-node/config"
)

type Signer interface {
	PubKey() cryptotypes.PubKey
	Sign(data []byte) ([]byte, error)
}

func NewSigner(ctx context.Context, config *configtypes.Config) Signer {
	if config.Feeder.AWSKMSKey != "" {
		return NewKmsSigner(ctx, config.Feeder.AWSKMSKey)
	}
	return NewLocalSigner(config.Feeder.PrivateKey)
}
