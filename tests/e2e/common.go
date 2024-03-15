package e2e

import "os"

var (
	chainAPIEndpoint = DefaultChainAPIEndpoint
	ethAPIEndpoint   = DefaultEthAPIEndpoint
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
