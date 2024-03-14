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

const masterPrivateKey = "f56cf1cd9f03c9a556da34ac8aefa1109c42d6e11d2e35ede699b515d0c7a56a"

func mintNFTContract(client *ethclient.Client) (string, error) {
	privateKey, err := crypto.HexToECDSA(masterPrivateKey)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(5371))
	if err != nil {
		return "", err
	}

	address, tx, _, err := bind.DeployContract(auth, contracts.ERC721Contract.ABI, contracts.ERC721Contract.Bin, client, "E2E_NFT", "E2E_NFT")
	if err != nil {
		fmt.Println("failed to deply contract", err)

		return "", err
	}
	_, err = bind.WaitDeployed(context.TODO(), client, tx)
	if err != nil {
		fmt.Println("failed to wait deployed", err)
	}

	return address.Hex(), nil
}

func mintNFT(client *ethclient.Client, contractAddr string, toAddress string) error {
	privateKey, err := crypto.HexToECDSA(masterPrivateKey)
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(5371))
	if err != nil {
		return err
	}

	contract := bind.NewBoundContract(common.HexToAddress(contractAddr), contracts.ERC721Contract.ABI, client, client, nil)
	tx, err := contract.Transact(auth, "safeMint", common.HexToAddress(toAddress))
	if err != nil {
		fmt.Println("failed to mint NFT")
		return err
	}

	_, err = bind.WaitMined(context.TODO(), client, tx)
	if err != nil {
		fmt.Println("failed to wait mined", err)
	}

	return nil
}

func queryERC20Balance(client *ethclient.Client, contractAddr string, addr string) (uint64, error) {
	contract := bind.NewBoundContract(common.HexToAddress(contractAddr), contracts.ERC20Contract.ABI, client, nil, nil)

	var result []interface{}
	balance := new(big.Int)
	result = append(result, &balance)
	err := contract.Call(nil, &result, "balanceOf", common.HexToAddress(addr))
	if err != nil {
		return 0, err
	}

	return balance.Uint64(), nil
}
