package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

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
	params := k.GetParams(ctx)
	blockHeight := ctx.BlockHeight()
	prevoteEnd, voteEnd := types.CalculateVotePeriod(blockHeight, params.VotePeriod)

	oracleData := []*types.OracleData{}
	oracleData = appendIfValid(oracleData, k.blockOracleData(params.Whitelist))
	oracleData = appendIfValid(oracleData, k.ownershipOracleData(ctx, params.VotePeriod))

	roundInfo := types.RoundInfo{
		Id:         types.CalculateRoundId(blockHeight, params.VotePeriod),
		PrevoteEnd: prevoteEnd,
		VoteEnd:    voteEnd,
		OracleData: oracleData,
		Timestamp:  ctx.BlockHeader().Time.Add(-BlockTimestampMargin).UnixMilli(),
	}

	return &roundInfo
}

func (k Keeper) blockOracleData(whitelist []*types.Chain) *types.OracleData {
	source := make([]string, len(whitelist))
	for idx, chain := range whitelist {
		source[idx] = chain.ChainId
	}

	return &types.OracleData{
		Topic:   types.OralceTopic_BLOCK,
		Sources: source,
	}
}

func (k Keeper) ownershipOracleData(ctx sdk.Context, votePeriod uint64) *types.OracleData {
	// We will retrieve all NFTs that need to be verified collected until the last round.
	// Because list of nfts in the current round can be increased as the round progresses
	startHeight := types.CalculateRoundStartHeight(ctx.BlockHeight(), votePeriod)
	nfts := k.SettlementKeeper.GetAllUniqueNftToVerify(ctx, startHeight-1)
	sources := make([]string, len(nfts))
	for i, nft := range nfts {
		sources[i] = nft.FormatString()
	}

	return &types.OracleData{
		Topic:   types.OralceTopic_OWNERSHIP,
		Sources: sources,
	}
}

// GetBlockData returns block data of a chain
func (k Keeper) GetBlockData(ctx sdk.Context, chainId string) (*types.BlockData, error) {
	store := ctx.KVStore(k.storeKey)
	chain, err := k.GetChain(ctx, chainId)
	if err != nil {
		return nil, fmt.Errorf("chain not found: %w", err)
	}

	bd := store.Get(types.BlockDataKey(chain.ChainId))
	if bd == nil {
		return nil, types.ErrBlockDataNotFound
	}
	var blockData types.BlockData
	k.cdc.MustUnmarshal(bd, &blockData)
	return &blockData, nil
}

// GetAllBlockData returns block data of all chains
func (k Keeper) GetAllBlockData(ctx sdk.Context) []types.BlockData {
	store := ctx.KVStore(k.storeKey)
	var blockData []types.BlockData
	iterator := sdk.KVStorePrefixIterator(store, types.BlockDataKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var bd types.BlockData
		k.cdc.MustUnmarshal(iterator.Value(), &bd)
		blockData = append(blockData, bd)
	}
	return blockData
}

// SetBlockData sets block data of a chain
func (k Keeper) SetBlockData(ctx sdk.Context, blockData types.BlockData) {
	store := ctx.KVStore(k.storeKey)
	bd := k.cdc.MustMarshal(&blockData)
	store.Set(types.BlockDataKey(blockData.ChainId), bd)
}

func (k Keeper) FillSettlementRecipients(ctx sdk.Context, nftOwnership map[types.Nft]ctypes.HexAddressString) {
	startHeight := types.CalculateRoundStartHeight(ctx.BlockHeight(), k.GetParams(ctx).VotePeriod)
	k.SettlementKeeper.SetRecipients(ctx, nftOwnership, startHeight-1)
}

func appendIfValid(slice []*types.OracleData, elem *types.OracleData) []*types.OracleData {
	if len(elem.Sources) == 0 {
		return slice
	}

	return append(slice, elem)
}
