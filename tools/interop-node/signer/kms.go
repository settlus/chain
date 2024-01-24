package signer

import (
	"context"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/settlus/chain/evmos/crypto/ethsecp256k1"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

type kmsSigner struct {
	ctx   context.Context
	svc   *kms.Client
	keyId string
	pub   cryptotypes.PubKey
}

func NewKmsSigner(ctx context.Context, key string) Signer {
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	svc := kms.NewFromConfig(awsCfg)

	signer := &kmsSigner{
		ctx:   ctx,
		svc:   svc,
		keyId: key,
	}

	if err := signer.loadPubKey(); err != nil {
		panic(fmt.Errorf("failed to load public key during kms signer initialization: %w", err))
	}

	return signer
}

func (s *kmsSigner) PubKey() cryptotypes.PubKey {
	return s.pub
}

func (s *kmsSigner) Sign(data []byte) ([]byte, error) {
	hash := crypto.Keccak256(data)
	output, err := s.svc.Sign(s.ctx, &kms.SignInput{
		KeyId:            aws.String(s.keyId),
		Message:          hash,
		MessageType:      types.MessageTypeDigest,
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
	})
	if err != nil {
		return nil, err
	}

	type ecdsaSignature struct {
		R, S *big.Int
	}

	var signature ecdsaSignature
	_, err = asn1.Unmarshal(output.Signature, &signature)
	if err != nil {
		return nil, err
	}

	if signature.S.Cmp(secp256k1halfN) > 0 {
		signature.S = signature.S.Sub(secp256k1N, signature.S)
	}

	rBytes := signature.R.Bytes()
	sBytes := signature.S.Bytes()
	sig := make([]byte, 64)

	copy(sig[32-len(rBytes):], rBytes)
	copy(sig[64-len(sBytes):], sBytes)
	return sig, nil
}

func (s *kmsSigner) loadPubKey() error {
	output, err := s.svc.GetPublicKey(s.ctx, &kms.GetPublicKeyInput{
		KeyId: aws.String(s.keyId),
	})

	if err != nil {
		return err
	}

	type publicKeyInfo struct {
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}

	var pki publicKeyInfo
	_, err = asn1.Unmarshal(output.PublicKey, &pki)
	if err != nil {
		return err
	}

	s256Pub, err := crypto.UnmarshalPubkey(pki.PublicKey.Bytes)
	if err != nil {
		return err
	}

	s.pub = &ethsecp256k1.PubKey{
		Key: crypto.CompressPubkey(s256Pub),
	}

	return nil
}
