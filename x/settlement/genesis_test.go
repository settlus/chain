package settlement_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	feemarkettypes "github.com/settlus/chain/evmos/x/feemarket/types"
	"github.com/settlus/chain/testutil/sample"

	"github.com/settlus/chain/app"
	utiltx "github.com/settlus/chain/testutil/tx"
	"github.com/settlus/chain/x/settlement"
	"github.com/settlus/chain/x/settlement/types"
)

type GenesisTestSuite struct {
	suite.Suite
	ctx     sdk.Context
	app     *app.App
	genesis types.GenesisState
	admin   string
	creator string
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	// consensus key
	consAddress := sdk.ConsAddress(utiltx.GenerateAddress().Bytes())
	checkTx := false

	suite.app = app.Setup(checkTx, feemarkettypes.DefaultGenesisState())
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
	suite.admin = sample.AccAddress()
	suite.creator = sample.AccAddress()
}

func (suite *GenesisTestSuite) TestGenesis_ExportGenesis() {
	testCases := []struct {
		name    string
		genesis types.GenesisState
	}{
		{
			name:    "default genesis",
			genesis: *types.DefaultGenesis(),
		},
		{
			name: "genesis with utxrs and tenants",
			genesis: types.GenesisState{
				Params: types.Params{
					GasPrice:            sdk.NewCoin("uusdc", sdk.NewInt(100)),
					OracleFeePercentage: sdk.NewDecWithPrec(1, 2),
				},
				Tenants: []types.Tenant{
					{
						Id:           1,
						Admins:       []string{sample.AccAddress()},
						PayoutPeriod: 100,
					},
				},
				Utxrs: []types.UTXRWithTenantAndId{
					{
						Id:       1,
						TenantId: 1,
						Utxr: types.UTXR{
							RequestId:  "request-0",
							Recipients: types.SingleRecipients(sdk.MustAccAddressFromBech32(suite.creator)),
							Amount:     sdk.NewCoin("uusdc", sdk.NewInt(100)),
							CreatedAt:  101,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			settlement.InitGenesis(suite.ctx, suite.app.SettlementKeeper, tc.genesis)
			genesis := settlement.ExportGenesis(suite.ctx, suite.app.SettlementKeeper)
			suite.Require().Equal(tc.genesis, *genesis)
		})
	}
}
