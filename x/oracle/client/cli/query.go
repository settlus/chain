package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/settlus/chain/x/oracle/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group oracle queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdQueryAggregatePrevote())
	cmd.AddCommand(CmdQueryAggregatePrevotes())
	cmd.AddCommand(CmdQueryAggregateVote())
	cmd.AddCommand(CmdQueryAggregateVotes())
	cmd.AddCommand(CmdQueryFeederDelegation())
	cmd.AddCommand(CmdQueryMissCount())
	cmd.AddCommand(CmdQueryRewardPool())
	cmd.AddCommand(CmdQueryCurrentRoundInfo())

	// this line is used by starport scaffolding # 1

	return cmd
}

// CmdQueryParams queries the parameters of oracle module
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

// CmdQueryAggregatePrevote queries an aggregate prevote of a validator.
func CmdQueryAggregatePrevote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-prevote [validator-address]",
		Short: "Query aggregate prevote",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqValidatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAggregatePrevoteRequest{
				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.AggregatePrevote(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryAggregatePrevotes queries a list of aggregate prevotes of all aggregate prevotes.
func CmdQueryAggregatePrevotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-prevotes",
		Short: "Query aggregate prevotes of all validators",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAggregatePrevotesRequest{Pagination: pageReq}

			res, err := queryClient.AggregatePrevotes(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryAggregateVote queries an aggregate vote of a validator.
func CmdQueryAggregateVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-vote [validator-address]",
		Short: "Query aggregate vote",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqValidatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAggregateVoteRequest{
				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.AggregateVote(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryAggregateVotes queries a list of aggregate votes of all aggregate votes.
func CmdQueryAggregateVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-votes",
		Short: "Query aggregate votes of all validators",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAggregateVotesRequest{Pagination: pageReq}

			res, err := queryClient.AggregateVotes(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryFeederDelegation queries a feeder delegation of a validator.
func CmdQueryFeederDelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeder-delegation [validator-address]",
		Short: "Query feeder delegation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqValidatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryFeederDelegationRequest{
				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.FeederDelegation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryMissCount queries a miss count of a validator.
func CmdQueryMissCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miss-count [validator-address]",
		Short: "Query miss count",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqValidatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryMissCountRequest{
				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.MissCount(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryRewardPool queries the current oracle reward pool balance.
func CmdQueryRewardPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-pool",
		Short: "Query reward pool",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRewardPoolRequest{}

			res, err := queryClient.RewardPool(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryCurrentRoundInfo queries the current oracle reward pool balance.
func CmdQueryCurrentRoundInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "round-info",
		Short: "Query current round info",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryCurrentRoundInfoRequest{}

			res, err := queryClient.CurrentRoundInfo(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
