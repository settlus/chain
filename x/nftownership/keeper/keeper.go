package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/settlus/chain/contracts"
	"github.com/settlus/chain/x/nftownership/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	evmKeeper     types.EVMKeeper
	oracleKeeper  types.OracleKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	accountKeeper types.AccountKeeper,
	evmKeeper types.EVMKeeper,
	oracleKeeper types.OracleKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,

		accountKeeper: accountKeeper,
		evmKeeper:     evmKeeper,
		oracleKeeper:  oracleKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) CheckValidChainId(ctx sdk.Context, chainId string) bool {
	params := k.GetParams(ctx)
	for _, allowedChainId := range params.AllowedChainIds {
		if allowedChainId == chainId {
			return true
		}
	}
	return false
}

func (k Keeper) OwnerOf(ctx sdk.Context, chainId string, contractAddr string, tokenIdHex string) (*common.Address, error) {
	if ctx.ChainID() == chainId {
		address, err := k.FindInternalOwner(ctx, contractAddr, tokenIdHex)
		if err != nil {
			return nil, types.ErrEVMCallFailed
		}
		return address, nil
	} else if k.CheckValidChainId(ctx, chainId) {
		return k.FindExternalOwner(ctx, chainId, contractAddr, tokenIdHex)
	} else {
		return &common.Address{}, types.ErrInvalidChainId
	}
}

// FindInternalOwner returns the owner of the given NFT on current chain
func (k Keeper) FindInternalOwner(
	ctx sdk.Context, contractAddr string, tokenIdHex string,
) (*common.Address, error) {
	erc721 := contracts.ERC721Contract.ABI
	tokenId := common.HexToHash(tokenIdHex)
	contract := common.HexToAddress(contractAddr)
	res, err := k.evmKeeper.CallEVM(ctx, erc721, types.ModuleAddress, contract, false, "ownerOf", tokenId.Big())
	if err != nil {
		return nil, fmt.Errorf("call evm failed: %w", err)
	}

	owner := common.BytesToAddress(res.Ret)
	if owner == common.HexToAddress("") {
		return nil, errorsmod.Wrapf(sdkerrorstypes.ErrInvalidAddress, "contract '%s', token '%s', owner '%s'", contract, tokenId, owner.String())
	}

	return &owner, nil
}

func (k Keeper) FindExternalOwner(ctx sdk.Context, chainId string, contractAddr string, tokenIdHex string) (*common.Address, error) {
	// Not implemented yet
	return nil, errortypes.ErrNotSupported
}
