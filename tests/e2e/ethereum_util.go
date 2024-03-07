package e2e

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/settlus/chain/contracts"

	"github.com/ethereum/go-ethereum/crypto"
)

func mintNFTContract(client *ethclient.Client) (string, error) {
	privateKey, err := crypto.HexToECDSA(adminPrivateKey)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(5371))
	if err != nil {
		return "", err
	}

	// ERC721 컨트랙트 배포
	address, tx, bound, err := bind.DeployContract(auth, contracts.ERC721Contract.ABI, contracts.ERC721Contract.Bin, client, "E2E_NFT", "E2E_NFT")
	if err != nil {
		fmt.Println("failed to deply contract", err)

		return "", err
	}
	bind.WaitDeployed(context.TODO(), client, tx)

	toAddress := common.HexToAddress(internalNftOwner)
	tx, err = bound.Transact(auth, "safeMint", toAddress)
	if err != nil {
		fmt.Println("failed to mint NFT")
		return "", err
	}
	bind.WaitMined(context.TODO(), client, tx)

	return address.Hex(), nil
}
