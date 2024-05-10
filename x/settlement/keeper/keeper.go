package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/settlus/chain/contracts"
	ctypes "github.com/settlus/chain/types"
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
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	erc20k types.Erc20Keeper,
	evmk types.EvmKeeper,
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
	}
}

func (k SettlementKeeper) InitAccountModule(ctx sdk.Context) {
	baseAcc := authtypes.NewBaseAccountWithAddress(authtypes.NewModuleAddress(types.ModuleName))
	accountName := fmt.Sprintf("%s-module-account", types.ModuleName)
	acc := authtypes.NewModuleAccount(baseAcc, accountName)
	k.ak.SetModuleAccount(ctx, acc)
}

func (k SettlementKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//NFT

// GetRecipient returns the owner of the given NFT from x/nftownership module
func (k SettlementKeeper) GetRecipients(ctx sdk.Context, chainId string, contractAddr string, tokenIdHex string) ([]*types.Recipient, error) {
	recipients := make([]*types.Recipient, 0)
	if k.IsSupportedChain(ctx, chainId) && ctx.ChainID() != chainId {
		// length 0 recipients means that owner of NFT will be determined by voting
		return recipients, nil
	}

	if ctx.ChainID() != chainId {
		return nil, errorsmod.Wrapf(types.ErrInvalidChainId, "chain '%s' is not supported", chainId)
	}

	address, err := k.FindInternalOwner(ctx, contractAddr, tokenIdHex)
	if err != nil {
		if sdkerrorstypes.ErrInvalidAddress.Is(err) {
			return nil, err
		}

		// Note: we shouldn't return the EVM error log directly to the client
		return nil, errorsmod.Wrapf(types.ErrEVMCallFailed, "failed to get the owner of the given NFT")
	}

	// TODO: support multi-owner
	recipient := &types.Recipient{
		Address: ctypes.NewHexAddrFromBytes(address.Bytes()),
		Weight:  1,
	}

	return append(recipients, recipient), nil
}

// FindInternalOwner returns the owner of the given NFT on current chain
func (k SettlementKeeper) FindInternalOwner(
	ctx sdk.Context, contractAddr string, tokenIdHex string,
) (*common.Address, error) {
	erc721 := contracts.ERC721Contract.ABI
	tokenId := common.HexToHash(tokenIdHex)
	contract := common.HexToAddress(contractAddr)
	res, err := k.evmk.CallEVM(ctx, erc721, types.ModuleAddress, contract, false, "ownerOf", tokenId.Big())
	if err != nil {
		return nil, fmt.Errorf("call evm failed: %w", err)
	}

	owner := common.BytesToAddress(res.Ret)
	if owner == common.HexToAddress("") {
		return nil, errorsmod.Wrapf(sdkerrorstypes.ErrInvalidAddress, "contract '%s', token '%s', owner '%s'", contract, tokenId, owner.String())
	}

	return &owner, nil
}
