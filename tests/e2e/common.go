package e2e

import "os"

var (
	chainAPIEndpoint = DefaultChainAPIEndpoint
	ethAPIEndpoint   = DefaultEthAPIEndpoint

	adminPrivateKey  = "f56cf1cd9f03c9a556da34ac8aefa1109c42d6e11d2e35ede699b515d0c7a56a"
	internalNftOwner = "0xa7801b8115f3fe46ac55f8c0fdb5243726bdb66a"
)

const (
	DefaultChainAPIEndpoint = "http://localhost:1317"
	DefaultEthAPIEndpoint   = "http://localhost:8545"
)

func init() {
	envChainApiEndPoint := os.Getenv("E2E_CHAIN_API_ENDPOINT")
	if envChainApiEndPoint != "" {
		chainAPIEndpoint = envChainApiEndPoint
	}

	envEthApiEndPoint := os.Getenv("E2E_ETH_API_ENDPOINT")
	if envEthApiEndPoint != "" {
		ethAPIEndpoint = envEthApiEndPoint
	}
}
