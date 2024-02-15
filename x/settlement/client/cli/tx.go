package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/settlus/chain/x/settlement/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
	FlagContractAddress                   = "contract-address"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdRecord())
	cmd.AddCommand(CmdCancel())
	cmd.AddCommand(CmdDepositToTreasury())
	cmd.AddCommand(CmdCreateTenant())
	cmd.AddCommand(CmdCreateTenantWithMintableContract())
	cmd.AddCommand(CmdAddTenantAdmin())
	cmd.AddCommand(CmdRemoveTenantAdmin())

	return cmd
}

func CmdRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record [tenant-id] [request-id] [amount] [chain-id] [contract-address] [token-id-hex] [metadata]",
		Short: "Record a settlement record in the settlement module",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new Unspent Transaction Record (UTXR) in the 'x/settlement' module.

Each UTXR is identified by a unique UTXR ID and is associated with a tenant, specified by the tenant ID. 
The contract address and the token ID (in hex format) specify the NFT that generated the revenue to be recorded.
The owner of the NFT will receive the payment in the specified amount.

Examples:
$ %s tx settlement record 1234 request-0 10usdc settlus_5371-1 
	0x0000000000000000000000000000000000000001 0x1 --from mykey  # NFT on Settlus

$ %s tx settlement record 1234 request-0 10usdc 1 
	0x0000000000000000000000000000000000000001 0x1 --from mykey  # NFT on external chains like Ethereum
`,
				version.AppName, version.AppName,
			),
		),
		Args: cobra.RangeArgs(6, 8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTenantId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argRequestId := args[1]
			argAmount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}
			argChainId := args[3]
			argContractAddress := args[4]
			tokenIdHex := args[5]

			var argMetadata string
			if len(args) >= 7 {
				argMetadata = args[6]
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRecord(
				clientCtx.GetFromAddress().String(),
				argTenantId,
				argRequestId,
				argAmount,
				argChainId,
				argContractAddress,
				tokenIdHex,
				argMetadata,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCancel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel [tenant-id] [request-id]",
		Short: "cancel a payment in the settlement module",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel a payment in the settlement module.

Cancellation is only possible when the payment is not settled yet.
If cancellation is eligable, coresponding UTXR is simply deleted from the state.

Example:
$ %s tx settlement cancel 1234 tx-5678 --from mykey
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTenantId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argRequestId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancel(
				clientCtx.GetFromAddress().String(),
				argTenantId,
				argRequestId,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDepositToTreasury() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-to-treasury [tenant-id] [amount]",
		Short: "Broadcast message depositToTreasury",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit to a tenant's treasury account. Anyone can deposit to a tenant's treasury account.

Example:
$ %s tx settlement deposit-to-treasury 1234 100usdc --from mykey
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTenantId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argAmount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositToTreasury(
				clientCtx.GetFromAddress().String(),
				argTenantId,
				argAmount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreateTenant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-tenant [denom] [payout-period]",
		Short: "Create a new tenant",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new tenant. The sender of the transaction will be the initial admin of the tenant.

Example:
$ %s tx settlement create-tenant asetl 100 --from mykey
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]

			argPayoutPeriod, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateTenant(
				clientCtx.GetFromAddress().String(),
				argDenom,
				argPayoutPeriod,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreateTenantWithMintableContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-tenant-mc [denom] [payout-period]",
		Short: "Create a new tenant with mintable contract",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new tenant with mintable contract functionality. The sender of the transaction
will be the initial admin of the tenant. The optional 'contract-address' flag can be used to specify the address
of the mintable contract. If not specified, a new mintable contract will be deployed.

Examples:
$ %s tx settlement create-tenant-mc asetl 100 --from mykey
$ %s tx settlement create-tenant-mc asetl 100 --contract-address 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48 --from mykey
`,
				version.AppName,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]

			argPayoutPeriod, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argContractAddr := cmd.Flag(FlagContractAddress).Value.String()
			msg := types.NewMsgCreateTenantWithMintableContract(
				clientCtx.GetFromAddress().String(),
				argDenom,
				argPayoutPeriod,
				argContractAddr,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagContractAddress, "", "address of mintable contract")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdAddTenantAdmin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-tenant-admin [tenant-id] [admin]",
		Short: "Add a new admin to the tenant",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add a new admin to the tenant. The sender of the transaction must be an admin of the tenant.

Example:
$ %s tx settlement add-tenant-admin 1234 settlus1... --from mykey
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTenantId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argAdmin := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddTenantAdmin(
				clientCtx.GetFromAddress().String(),
				argTenantId,
				argAdmin,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdRemoveTenantAdmin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-tenant-admin [tenant-id] [admin]",
		Short: "Remove an admin from the tenant",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Remove an admin from the tenant. The sender of the transaction must be an admin of the tenant.
Cannot remove the last admin.

Example:
$ %s tx settlement remove-tenant-admin 1234 settlus1... --from mykey
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTenantId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argAdmin := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveTenantAdmin(
				clientCtx.GetFromAddress().String(),
				argTenantId,
				argAdmin,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
