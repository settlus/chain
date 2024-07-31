package oracle_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/tmhash"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	"github.com/cometbft/cometbft/version"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	feemarkettypes "github.com/evmos/evmos/v19/x/feemarket/types"

	"github.com/settlus/chain/app"
	utiltx "github.com/settlus/chain/testutil/tx"
	"github.com/settlus/chain/utils"
	"github.com/settlus/chain/x/oracle"
	"github.com/settlus/chain/x/oracle/types"
)

type GenesisTestSuite struct {
	suite.Suite
	ctx     sdk.Context
	app     *app.SettlusApp
	genesis types.GenesisState
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	// consensus key
	consAddress := sdk.ConsAddress(utiltx.GenerateAddress().Bytes())
	checkTx := false

	suite.app = app.Setup(checkTx, feemarkettypes.DefaultGenesisState(), utils.MainnetChainID)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{
		Height:          1,
		ChainID:         "settlus_5371-1",
		Time:            time.Now().UTC(),
		ProposerAddress: consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})

	suite.genesis = *types.DefaultGenesis()
}

func (suite *GenesisTestSuite) TestOracleInitGenesis() {
	testCases := []struct {
		name         string
		genesisState types.GenesisState
	}{
		{
			"default genesis",
			*types.DefaultGenesis(),
		},
	}

	for _, tc := range testCases {
		suite.Require().NotPanics(func() {
			oracle.InitGenesis(suite.ctx, *suite.app.OracleKeeper, tc.genesisState)
		})
		params := suite.app.OracleKeeper.GetParams(suite.ctx)

		suite.Require().Equal(tc.genesisState.Params, params)
	}
}
