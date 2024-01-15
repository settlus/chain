package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/settlus/chain/x/nftownership/types"
)

var _ = strconv.Itoa(0)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group nftownership queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdGetNFTOwner())

	// this line is used by starport scaffolding # 1

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetNFTOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-nft-owner [chain-id] [contract-address] [token-id-hex]",
		Short: "get nft owner",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			reqChainId := args[0]
			reqContractAddress := args[1]
			reqTokenId := args[2]

			contractAddress := common.HexToAddress(reqContractAddress)
			if contractAddress == common.HexToAddress("") {
				return fmt.Errorf("invalid contract address")
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetNftOwnerRequest{
				ChainId:         reqChainId,
				ContractAddress: contractAddress.String(),
				TokenIdHex:      reqTokenId,
			}

			res, err := queryClient.GetNftOwner(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
