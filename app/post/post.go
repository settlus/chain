package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmospost "github.com/evmos/evmos/v19/app/post"
)

func NewPostHandler(ho evmospost.HandlerOptions) sdk.PostHandler {
	postDecorators := []sdk.PostDecorator{
		NewSettlementDecorator(),
		evmospost.NewBurnDecorator(ho.FeeCollectorName, ho.BankKeeper),
	}

	return sdk.ChainPostDecorators(postDecorators...)
}
