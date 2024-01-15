package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/settlus/chain/x/settlement/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
)

type SettlementKeeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramstore paramtypes.Subspace
	ak         types.AccountKeeper
	bk         types.BankKeeper
	erc20k     types.Erc20Keeper
	evmk       types.EvmKeeper
	nk         types.NftOwnershipKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	erc20k types.Erc20Keeper,
	evmk types.EvmKeeper,
	nk types.NftOwnershipKeeper,
) *SettlementKeeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &SettlementKeeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,
		ak:         ak,
		bk:         bk,
		erc20k:     erc20k,
		evmk:       evmk,
		nk:         nk,
	}
}

func (k SettlementKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetRecipient returns the owner of the given NFT from x/nftownership module
func (k SettlementKeeper) GetRecipient(ctx sdk.Context, chaindId string, contractAddr string, tokenIdHex string) (string, error) {
	ownerAddress, err := k.nk.OwnerOf(ctx, chaindId, contractAddr, tokenIdHex)
	if err != nil {
		if sdkerrorstypes.ErrInvalidAddress.Is(err) {
			return "", err
		}

		// Note: we shouldn't return the EVM error log directly to the client
		return "", errorsmod.Wrapf(types.ErrEVMCallFailed, "failed to get the owner of the given NFT")
	}

	return ownerAddress.Hex(), nil
}
