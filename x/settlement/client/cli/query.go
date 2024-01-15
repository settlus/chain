package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/settlus/chain/x/settlement/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group settlement queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdQueryUTXR())
	cmd.AddCommand(CmdQueryUTXRAll())
	cmd.AddCommand(CmdTenant())
	cmd.AddCommand(CmdTenants())

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the x/settlement module",
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

func CmdQueryUTXR() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "utxr [tenant-id] [request-id]",
		Short: "shows the UTXR of the x/settlement module",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a UTXR by tenant ID and request ID.

Example:
$ %s query settlement utxr 1234 tx-5678
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tenantId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryUTXRRRequest{
				TenantId:  tenantId,
				RequestId: args[1],
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.UTXR(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryUTXRAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "utxrs [tenant-id]",
		Short: "shows the UTXR of the x/settlement module",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all UTXR by tenant ID.

Example:
$ %s query settlement utxrs 1234
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tenantId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryUTXRsRequest{TenantId: tenantId, Pagination: pageReq}

			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.UTXRs(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdTenant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenant [tenant-id]",
		Short: "shows the tenant details of the x/settlement module",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a tenant by its ID.

Example:
$ %s query settlement tenant 1234
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tenantId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryTenantRequest{TenantId: tenantId}

			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Tenant(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdTenants() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenants",
		Short: "shows the tenants of the x/settlement module",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all tenants.

Example:
$ %s query settlement tenants
`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryTenantsRequest{Pagination: pageReq}

			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Tenants(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
