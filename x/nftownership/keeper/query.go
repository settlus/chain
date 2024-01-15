package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/settlus/chain/x/nftownership/types"
)

func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	wctx := sdk.UnwrapSDKContext(ctx)

	return &types.QueryParamsResponse{Params: k.GetParams(wctx)}, nil
}

func (k Keeper) GetNftOwner(goCtx context.Context, req *types.QueryGetNftOwnerRequest) (*types.QueryGetNftOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	ownerAddress, err := k.OwnerOf(ctx, req.ChainId, req.ContractAddress, req.TokenIdHex)
	if err != nil {
		if sdkerrorstypes.ErrInvalidAddress.Is(err) {
			return nil, sdkerrorstypes.ErrInvalidAddress
		}

		// Note: we shouldn't return the EVM error log directly to the client
		return nil, err
	}
	return &types.QueryGetNftOwnerResponse{
		OwnerAddress: ownerAddress.String(),
	}, nil
}
