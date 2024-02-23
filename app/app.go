package app

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	appparams "github.com/cosmos/cosmos-sdk/simapp/params"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts"
	icahost "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibctransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v6/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v6/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v6/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"
	ibctestingtypes "github.com/cosmos/ibc-go/v6/testing/types"

	"github.com/settlus/chain/evmos/encoding"

	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	evmosante "github.com/settlus/chain/evmos/app/ante"
	ethante "github.com/settlus/chain/evmos/app/ante/evm"
	srvflags "github.com/settlus/chain/evmos/server/flags"
	evmostypes "github.com/settlus/chain/evmos/types"
	"github.com/settlus/chain/evmos/x/erc20"
	erc20client "github.com/settlus/chain/evmos/x/erc20/client"
	erc20keeper "github.com/settlus/chain/evmos/x/erc20/keeper"
	erc20types "github.com/settlus/chain/evmos/x/erc20/types"
	"github.com/settlus/chain/evmos/x/evm"
	evmkeeper "github.com/settlus/chain/evmos/x/evm/keeper"
	evmtypes "github.com/settlus/chain/evmos/x/evm/types"
	"github.com/settlus/chain/evmos/x/feemarket"
	feemarketkeeper "github.com/settlus/chain/evmos/x/feemarket/keeper"
	feemarkettypes "github.com/settlus/chain/evmos/x/feemarket/types"

	// NOTE: override ICS20 keeper to support IBC transfers of ERC20 tokens
	"github.com/settlus/chain/evmos/x/ibc/transfer"
	transferkeeper "github.com/settlus/chain/evmos/x/ibc/transfer/keeper"

	"github.com/settlus/chain/app/ante"

	"github.com/settlus/chain/swagger"
	"github.com/settlus/chain/x/interop"
	nftownershipmodule "github.com/settlus/chain/x/nftownership"
	nftownershipmodulekeeper "github.com/settlus/chain/x/nftownership/keeper"
	nftownershipmoduletypes "github.com/settlus/chain/x/nftownership/types"
	oraclemodule "github.com/settlus/chain/x/oracle"
	oraclemodulekeeper "github.com/settlus/chain/x/oracle/keeper"
	oraclemoduletypes "github.com/settlus/chain/x/oracle/types"
	settlementmodule "github.com/settlus/chain/x/settlement"
	settlementmodulekeeper "github.com/settlus/chain/x/settlement/keeper"
	settlementmoduletypes "github.com/settlus/chain/x/settlement/types"
	// this line is used by starport scaffolding # stargate/app/moduleImport
)

const (
	Name = "settlus"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.LegacyProposalHandler, upgradeclient.LegacyCancelProposalHandler,
				ibcclientclient.UpdateClientProposalHandler, ibcclientclient.UpgradeProposalHandler,
				// Evmos proposal types
				erc20client.RegisterCoinProposalHandler, erc20client.RegisterERC20ProposalHandler, erc20client.ToggleTokenConversionProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		groupmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{AppModuleBasic: &ibctransfer.AppModuleBasic{}},
		ica.AppModuleBasic{},
		// this line is used by starport scaffolding # stargate/app/moduleBasic

		// EVMOS app modules
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
		erc20.AppModuleBasic{},

		// Settlus modules
		settlementmodule.AppModuleBasic{},
		nftownershipmodule.AppModuleBasic{},
		oraclemodule.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          {authtypes.Burner},
		icatypes.ModuleName:            nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		// this line is used by starport scaffolding # stargate/app/maccPerms
		evmtypes.ModuleName:   {authtypes.Minter, authtypes.Burner}, // used for secure addition and subtraction of balance using module account
		erc20types.ModuleName: {authtypes.Minter, authtypes.Burner},

		// settlus modules
		settlementmoduletypes.ModuleName:   nil,
		nftownershipmoduletypes.ModuleName: nil,
		oraclemoduletypes.ModuleName:       nil,
	}
)

var (
	_ servertypes.Application = (*App)(nil)
	_ ibctesting.TestingApp   = (*App)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)
	sdk.ConstantReward = true
}

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	AuthzKeeper      authzkeeper.Keeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   transferkeeper.Keeper
	ICAHostKeeper    icahostkeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	GroupKeeper      groupkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper

	// this line is used by starport scaffolding # stargate/app/keeperDeclaration
	// Evmos keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper
	Erc20Keeper     erc20keeper.Keeper

	// Settlus keepers
	SettlementKeeper   *settlementmodulekeeper.SettlementKeeper
	NftOwnershipKeeper *nftownershipmodulekeeper.Keeper
	OracleKeeper       *oraclemodulekeeper.Keeper

	// mm is the module manager
	mm *module.Manager

	// the configurator
	configurator module.Configurator

	// sm is the simulation manager
	sm *module.SimulationManager

	tpsCounter *tpsCounter
}

// New returns a reference to an initialized blockchain app
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig appparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, capabilitytypes.StoreKey,
		feegrant.StoreKey, authzkeeper.StoreKey,
		group.StoreKey, minttypes.StoreKey,
		// this line is used by starport scaffolding # stargate/app/storeKey
		// ibc keys
		ibchost.StoreKey, ibctransfertypes.StoreKey,
		// ica keys
		icahosttypes.StoreKey,
		// Ethermint keys
		evmtypes.StoreKey, feemarkettypes.StoreKey,
		// Evmos keys
		erc20types.StoreKey,
		// Settlus keys
		settlementmoduletypes.StoreKey,
		nftownershipmoduletypes.StoreKey,
		oraclemoduletypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(
		appCodec,
		cdc,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	// this line is used by starport scaffolding # stargate/app/scopedKeeper

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		evmostypes.ProtoAccount,
		maccPerms,
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.BlockedAddrs(),
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		keys[minttypes.StoreKey],
		app.GetSubspace(minttypes.ModuleName),
		&stakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.GetSubspace(distrtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)

	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)

	groupConfig := group.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	app.GroupKeeper = groupkeeper.NewKeeper(
		keys[group.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		groupConfig,
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// ... other modules keepers
	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))

	// Create Ethermint keepers
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec, authtypes.NewModuleAddress(govtypes.ModuleName),
		keys[feemarkettypes.StoreKey],
		tkeys[feemarkettypes.TransientKey],
		app.GetSubspace(feemarkettypes.ModuleName),
	)

	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec, keys[evmtypes.StoreKey], tkeys[evmtypes.TransientKey], authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, &stakingKeeper, app.FeeMarketKeeper,
		tracer, app.GetSubspace(evmtypes.ModuleName),
	)

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), &stakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(erc20types.RouterKey, erc20.NewErc20ProposalHandler(&app.Erc20Keeper))

	govConfig := govtypes.DefaultConfig()
	/*
		Example of setting gov params:
		govConfig.MaxMetadataLen = 10000
	*/
	app.GovKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter, app.MsgServiceRouter(), govConfig,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	// NOTE: Distr, Slashing and Claim must be created before calling the Hooks method to avoid returning a Keeper without its table generated
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
		),
	)

	app.Erc20Keeper = erc20keeper.NewKeeper(
		keys[erc20types.StoreKey], appCodec, authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.EvmKeeper, app.StakingKeeper,
	)

	app.EvmKeeper = app.EvmKeeper.SetHooks(
		evmkeeper.NewMultiEvmHooks(
			app.Erc20Keeper.Hooks(),
		),
	)

	// Settlus Keeper
	app.OracleKeeper = oraclemodulekeeper.NewKeeper(
		appCodec,
		keys[oraclemoduletypes.StoreKey],
		app.GetSubspace(oraclemoduletypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		app.StakingKeeper,
		distrtypes.ModuleName,
	)

	interopNodePort := cast.ToUint16(appOpts.Get(interop.FlagInteropNodePort))
	interopClientFactory := interop.NewInteropClientFactory(logger, interopNodePort)
	app.NftOwnershipKeeper = nftownershipmodulekeeper.NewKeeper(
		appCodec,
		keys[nftownershipmoduletypes.StoreKey],
		app.GetSubspace(nftownershipmoduletypes.ModuleName),
		app.AccountKeeper,
		app.EvmKeeper,
		app.OracleKeeper,
		interopClientFactory,
	)

	app.SettlementKeeper = settlementmodulekeeper.NewKeeper(
		appCodec,
		keys[settlementmoduletypes.StoreKey],
		app.GetSubspace(settlementmoduletypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.Erc20Keeper,
		app.EvmKeeper,
		app.NftOwnershipKeeper,
	)

	// NOTE: app.Erc20Keeper is already initialized elsewhere

	app.TransferKeeper = transferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
		app.Erc20Keeper, // Add ERC20 Keeper for ERC20 transfers
	)

	// Create the app.ICAHostKeeper
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, app.keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)

	// create host IBC module
	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	/*
		Create Transfer Stack

		transfer stack contains (from bottom to top):
			- ERC-20 Middleware
			- IBC Transfer

		SendPacket, since it is originating from the application to core IBC:
		 	transferKeeper.SendPacket-> erc20.SendPacket -> channel.SendPacket

		RecvPacket, message that originates from core IBC and goes down to app, the flow is the other way
			channel.RecvPacket -> erc20.OnRecvPacket -> transfer.OnRecvPacket
	*/

	// create IBC module from top to bottom of stack
	var transferStack ibcporttypes.IBCModule

	transferStack = transfer.NewIBCModule(app.TransferKeeper)
	transferStack = erc20.NewIBCMiddleware(app.Erc20Keeper, transferStack)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibctransfertypes.ModuleName, transferStack)

	// this line is used by starport scaffolding # ibc/app/router
	app.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, minttypes.DefaultInflationCalculationFn),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		// this line is used by starport scaffolding # stargate/app/appModule
		// ibc modules
		ibc.NewAppModule(app.IBCKeeper),
		ica.NewAppModule(nil, &app.ICAHostKeeper),
		transfer.NewAppModule(app.TransferKeeper),

		// Evmos app modules
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.GetSubspace(evmtypes.ModuleName)),
		feemarket.NewAppModule(app.FeeMarketKeeper, app.GetSubspace(feemarkettypes.ModuleName)),
		erc20.NewAppModule(app.Erc20Keeper, app.AccountKeeper, app.GetSubspace(erc20types.ModuleName)),

		// Settlus app modules
		settlementmodule.NewAppModule(appCodec, app.SettlementKeeper, app.AccountKeeper, app.BankKeeper),
		nftownershipmodule.NewAppModule(appCodec, *app.NftOwnershipKeeper, app.AccountKeeper, app.EvmKeeper),
		oraclemodule.NewAppModule(appCodec, *app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		// upgrades should be run first
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/beginBlockers
		// Evmos modules
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		erc20types.ModuleName,

		// Settlus modules
		settlementmoduletypes.ModuleName,
		nftownershipmoduletypes.ModuleName,
		oraclemoduletypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/endBlockers
		// Evmos modules
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		erc20types.ModuleName,

		// Settlus modules
		settlementmoduletypes.ModuleName,
		nftownershipmoduletypes.ModuleName,
		oraclemoduletypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/initGenesis
		// Evmos modules
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		erc20types.ModuleName,
		// NOTE: feemarket module needs to be initialized before genutil module:
		// gentx transactions use MinGasPriceDecorator.AnteHandle
		genutiltypes.ModuleName,

		// Settlus modules
		settlementmoduletypes.ModuleName,
		nftownershipmoduletypes.ModuleName,
		oraclemoduletypes.ModuleName,

		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	anteHandler := ante.NewAnteHandler(ante.HandlerOptions{
		Cdc:                    app.appCodec,
		AccountKeeper:          app.AccountKeeper,
		BankKeeper:             app.BankKeeper,
		ExtensionOptionChecker: evmostypes.HasDynamicFeeExtensionOption,
		EvmKeeper:              app.EvmKeeper,
		StakingKeeper:          app.StakingKeeper,
		IBCKeeper:              app.IBCKeeper,
		FeeMarketKeeper:        app.FeeMarketKeeper,
		FeegrantKeeper:         app.FeeGrantKeeper,
		SignModeHandler:        encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:         evmosante.SigVerificationGasConsumer,
		MaxTxGasWanted:         cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted)),
		SettlementKeeper:       app.SettlementKeeper,
		TxFeeChecker:           ethante.NewDynamicFeeChecker(app.EvmKeeper),
	})

	app.SetAnteHandler(anteHandler)
	app.SetPostHandler(ante.NewPostHandler(ante.HandlerOptions{}))
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	// this line is used by starport scaffolding # stargate/app/beforeInitReturn
	// Finally start the tpsCounter.
	app.tpsCounter = newTPSCounter(logger)
	go func() {
		// Unfortunately golangci-lint is so pedantic
		// so we have to ignore this error explicitly.
		_ = app.tpsCounter.start(context.Background())
	}()

	return app
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// BeginBlocker runs the Tendermint ABCI BeginBlock logic. It executes state changes at the beginning
// of the new block for every registered module. If there is a registered fork at the current height,
// BeginBlocker will schedule the upgrade plan and perform the state migration (if any).
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// The DeliverTx method is intentionally decomposed to calculate the transactions per second.
func (app *App) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	defer func() {
		// TODO: Record the count along with the code and or reason so as to display
		// in the transactions per second live dashboards.
		if res.IsErr() {
			app.tpsCounter.incrementFailure()
		} else {
			app.tpsCounter.incrementSuccess()
		}
	}()
	return app.BaseApp.DeliverTx(req)
}

// InitChainer updates at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads state at a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)

	accs := make([]string, 0, len(maccPerms))
	for k := range maccPerms {
		accs = append(accs, k)
	}
	sort.Strings(accs)

	for _, acc := range accs {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *App) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)

	accs := make([]string, 0, len(maccPerms))
	for k := range maccPerms {
		accs = append(accs, k)
	}
	sort.Strings(accs)

	for _, acc := range accs {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = true
		// blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc] 안씀
	}

	return blockedAddrs
}

// LegacyAmino returns Evmos's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns Evmos's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Evmos's InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	node.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register app's OpenAPI routes.
	if apiConfig.Swagger {
		swagger.RegisterOpenAPIService(Name, apiSvr.Router)
	}
}

func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterNodeService registers the node gRPC service on the provided
// application gRPC query router.
func (app *App) RegisterNodeService(clientCtx client.Context) {
	node.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// IBC Go TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *App) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetStakingKeeper implements the TestingApp interface.
func (app *App) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return app.StakingKeeper
}

// GetStakingKeeperSDK implements the TestingApp interface.
func (app *App) GetStakingKeeperSDK() stakingkeeper.Keeper {
	return app.StakingKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *App) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *App) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

// GetTxConfig implements the TestingApp interface.
func (app *App) GetTxConfig() client.TxConfig {
	cfg := encoding.MakeConfig(ModuleBasics)
	return cfg.TxConfig
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	// this line is used by starport scaffolding # stargate/app/paramSubspace
	// Evmos modules
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	paramsKeeper.Subspace(erc20types.ModuleName)

	// Settlus modules
	paramsKeeper.Subspace(settlementmoduletypes.ModuleName)
	paramsKeeper.Subspace(nftownershipmoduletypes.ModuleName)
	paramsKeeper.Subspace(oraclemoduletypes.ModuleName)

	return paramsKeeper
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}
