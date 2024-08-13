package keeper

import (
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	ctypes "github.com/settlus/chain/types"
	"github.com/settlus/chain/x/oracle/types"
)

const BlockTimestampMargin = 15 * time.Second

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace

		AccountKeeper      types.AccountKeeper
		BankKeeper         types.BankKeeper
		DistributionKeeper types.DistributionKeeper
		StakingKeeper      types.StakingKeeper
		SettlementKeeper   types.SettlementKeeper

		distributionName string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distributionKeeper types.DistributionKeeper,
	stakingKeeper types.StakingKeeper,
	settlementKeeper types.SettlementKeeper,

	distributionName string,
) *Keeper {
	// ensure oracle module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,

		AccountKeeper:      accountKeeper,
		BankKeeper:         bankKeeper,
		DistributionKeeper: distributionKeeper,
		StakingKeeper:      stakingKeeper,
		SettlementKeeper:   settlementKeeper,

		distributionName: distributionName,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetCurrentRoundInfo(ctx sdk.Context) *types.RoundInfo {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RoundKeyPrefix)
	if len(bz) == 0 {
		return nil
	}

	var roundInfo types.RoundInfo
	k.cdc.MustUnmarshal(bz, &roundInfo)
	return &roundInfo
}

func (k Keeper) SetCurrentRoundInfo(ctx sdk.Context, roundInfo *types.RoundInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(roundInfo)
	store.Set(types.RoundKeyPrefix, bz)
}

func (k Keeper) CalculateNextRoundInfo(ctx sdk.Context) *types.RoundInfo {
	params := k.GetParams(ctx)
	blockHeight := ctx.BlockHeight() + 1
	prevoteEnd, voteEnd := types.CalculateVotePeriod(blockHeight, params.VotePeriod)

	roundInfo := types.RoundInfo{
		Id:         types.CalculateRoundId(blockHeight, params.VotePeriod),
		PrevoteEnd: prevoteEnd,
		VoteEnd:    voteEnd,
		Ownerships: k.ownershipOracleData(ctx, params.VotePeriod),
		Timestamp:  ctx.BlockHeader().Time.Add(-BlockTimestampMargin).UnixMilli(),
	}

	return &roundInfo
}

func (k Keeper) ownershipOracleData(ctx sdk.Context, votePeriod uint64) []string {
	// We will retrieve all NFTs that need to be verified collected until the last round.
	// Because list of nfts in the current round can be increased as the round progresses
	startHeight := types.CalculateRoundStartHeight(ctx.BlockHeight(), votePeriod)
	nfts := k.SettlementKeeper.GetAllUniqueNftToVerify(ctx, startHeight-1)
	sources := make([]string, len(nfts))
	for i, nft := range nfts {
		sources[i] = nft.FormatString()
	}

	return sources
}

func (k Keeper) FillSettlementRecipients(ctx sdk.Context, nftOwnership map[ctypes.Nft]ctypes.HexAddressString) {
	startHeight := types.CalculateRoundStartHeight(ctx.BlockHeight(), k.GetParams(ctx).VotePeriod)
	k.SettlementKeeper.SetRecipients(ctx, nftOwnership, startHeight-1)
}
