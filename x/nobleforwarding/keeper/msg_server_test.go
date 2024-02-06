package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	keepertest "github.com/settlus/chain/testutil/keeper"
	"github.com/settlus/chain/x/nobleforwarding/keeper"
	"github.com/settlus/chain/x/nobleforwarding/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.NobleForwardingKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
