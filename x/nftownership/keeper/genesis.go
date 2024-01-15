package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/settlus/chain/x/nftownership/types"
)

/*
InitGenesis initializes the module's state from a provided genesis state.
If we do not set the module account here, we cannot query the EVM state with nftownership's module account.
*/
func (k Keeper) InitGenesis(ctx sdk.Context) {
	baseAcc := authtypes.NewBaseAccountWithAddress(authtypes.NewModuleAddress(types.ModuleName))
	accountName := fmt.Sprintf("%s-module-account", types.ModuleName)
	acc := authtypes.NewModuleAccount(baseAcc, accountName)
	k.accountKeeper.SetModuleAccount(ctx, acc)
}
