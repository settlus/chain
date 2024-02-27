package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/settlus/chain/x/oracle/types"
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

	cmd.AddCommand(CmdPrevote())
	cmd.AddCommand(CmdVote())
	cmd.AddCommand(CmdFeederDelegationConsent())
	// this line is used by starport scaffolding # 1

	return cmd
}

// CmdPrevote is the CLI command for sending a Prevote message
func CmdPrevote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prevote [validator] [hash] [roundId]",
		Short: "Broadcast message prevote",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argValidator := args[0]
			argHash := args[1]
			argRoundId := args[2]
			roundId, err := cast.ToUint64E(argRoundId)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			feeder := clientCtx.GetFromAddress().String()

			msg := types.NewMsgPrevote(
				feeder,
				argValidator,
				argHash,
				roundId,
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

// CmdVote is the CLI command for sending a Vote message
func CmdVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [validator] [blockData] [salt] [roundId]",
		Short: "Broadcast message vote",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argValidator := args[0]
			argBlockDataStr := args[1]
			argSalt := args[2]
			argRoundId := args[3]
			roundId, err := cast.ToUint64E(argRoundId)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			feeder := clientCtx.GetFromAddress().String()

			msg := types.NewMsgVote(
				feeder,
				argValidator,
				argBlockDataStr,
				argSalt,
				roundId,
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

// CmdFeederDelegationConsent is the CLI command for sending a FeederDelegationConsent message
func CmdFeederDelegationConsent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeder-delegation-consent [feeder-address]",
		Short: "Broadcast message feeder-delegation-consent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argFeederAddress := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			voter := sdk.ValAddress(clientCtx.GetFromAddress()).String()

			msg := types.NewMsgFeederDelegationConsent(
				voter,
				argFeederAddress,
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
