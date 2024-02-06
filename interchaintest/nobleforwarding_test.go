package interchaintest_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	chantypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/relayer"
	"github.com/strangelove-ventures/interchaintest/v6/testreporter"
	"github.com/strangelove-ventures/interchaintest/v6/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func Test_NobleForwardingAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Parallel()

	client, network := interchaintest.DockerSetup(t)

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()

	// Get both chains
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name: "settlusd",
			ChainConfig: ibc.ChainConfig{ // TODO: this should use settlusd binary or docker
				Type: "cosmos",
			},
		},
		{
			Name: "gaiad",
			ChainConfig: ibc.ChainConfig{ // TODO: use nobled binary when john implements the automated ica
				Type: "cosmos",
			},
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	settlusChain, cosmosChain := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Get a relayer instance
	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
		relayer.RelayerOptionExtraStartFlags{Flags: []string{"-p", "events", "-b", "100"}},
	).Build(t, client, network)

	// Build the network; spin up the chains and configure the relayer
	const pathName = "test-path"
	const relayerName = "relayer"

	ic := interchaintest.NewInterchain().
		AddChain(settlusChain).
		AddChain(cosmosChain).
		AddRelayer(r, relayerName).
		AddLink(interchaintest.InterchainLink{
			Chain1:  settlusChain,
			Chain2:  cosmosChain,
			Relayer: r,
			Path:    pathName,
		})

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
	}))

	// Fund a user account on settlusChain and cosmosChain
	userFunds := math.NewIntFromUint64(uint64(10_000_000_000))
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, settlusChain, cosmosChain)
	settlusUser := users[0]
	cosmosUser := users[1]

	// Generate a new IBC path
	err = r.GeneratePath(ctx, eRep, settlusChain.Config().ChainID, cosmosChain.Config().ChainID, pathName)
	require.NoError(t, err)

	// Create new clients
	err = r.CreateClients(ctx, eRep, pathName, ibc.CreateClientOptions{TrustingPeriod: "330h"})
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, settlusChain, cosmosChain)
	require.NoError(t, err)

	// Create a new connection
	err = r.CreateConnections(ctx, eRep, pathName)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, settlusChain, cosmosChain)
	require.NoError(t, err)

	// Query for the newly created connection
	connections, err := r.GetConnections(ctx, eRep, settlusChain.Config().ChainID)
	require.NoError(t, err)
	require.Equal(t, 1, len(connections))

	// Start the relayer and set the cleanup function.
	err = r.StartRelayer(ctx, eRep, pathName)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	// TODO: fix
	settlusAddr := settlusUser(*cosmos.CosmosWallet).FormattedAddress(settlusChain.Config().Bech32Prefix)

	registerForwardingAccount := []string{
		settlusChain.Config().Bin, "tx", "nobleforwarding", "register",
		"--from", settlusAddr,
		"--connection-id", connections[0].ID,
		"--chain-id", settlusChain.Config().ChainID,
		"--home", settlusChain.HomeDir(),
		"--node", settlusChain.GetGRPCAddress(),
		"--keyring-backend", keyring.BackendTest,
		"-y",
	}

	_, _, err = settlusChain.Exec(ctx, registerForwardingAccount, nil)
	require.NoError(t, err)

	ir := cosmos.DefaultEncoding().InterfaceRegistry

	cHeight, err := cosmosChain.Height(ctx)
	require.NoError(t, err)

	channelFound := func(found *chantypes.MsgChannelOpenConfirm) bool {
		return true // TODO: implement this
	}

	// Wait for channel open confirm
	_, err = cosmos.PollForMessage(ctx, cosmosChain, ir,
		cHeight, cHeight+30, channelFound)
	require.NoError(t, err)
}
