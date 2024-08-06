package signer

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/evmos/evmos/v19/crypto/ethsecp256k1"
)

type localSigner struct {
	privKey cryptotypes.PrivKey
}

func NewLocalSigner(key string) Signer {
	ecdsaPriv, err := crypto.HexToECDSA(key)
	if err != nil {
		panic(err)
	}

	privKey := &ethsecp256k1.PrivKey{
		Key: crypto.FromECDSA(ecdsaPriv),
	}

	return &localSigner{privKey: privKey}
}

func (s *localSigner) PubKey() cryptotypes.PubKey {
	return s.privKey.PubKey()
}

func (s *localSigner) Sign(data []byte) ([]byte, error) {
	return s.privKey.Sign(data)
}
