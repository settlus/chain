package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	channelutils "github.com/cosmos/ibc-go/v6/modules/core/04-channel/client/utils"
	"github.com/spf13/cobra"

	"github.com/settlus/chain/x/nobleforwarding/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagTimeoutTimestamp = "timeout-timestamp"
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

	cmd.AddCommand(CmdRegister())

	return cmd
}

// CmdRegister is the CLI command for registering an IBC forwarding account
func CmdRegister() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [src-port] [src-channel]",
		Short: "Register an IBC forwarding account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register an IBC forwarding account on the specified source port and channel.

The timeout timestamp is the timestamp at which the account will be closed if no packets are sent through it.
If the timeout timestamp is not provided, the default relative packet timeout timestamp will be used.

Example:
$ %s tx nobleforwarding register transfer channel-0 --timeout-timestamp 1000
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argPort := args[0]
			argChannel := args[1]
			argTimeoutTimestamp, err := cmd.Flags().GetUint64(flagTimeoutTimestamp)
			if err != nil {
				return fmt.Errorf("invalid timeout timestamp: %w", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			consensusState, _, _, err := channelutils.QueryLatestConsensusState(clientCtx, argPort, argChannel)
			if err != nil {
				return err
			}

			timeoutTimestamp := consensusState.GetTimestamp()
			if argTimeoutTimestamp != 0 {
				timeoutTimestamp += timeoutTimestamp
			} else {
				timeoutTimestamp += DefaultRelativePacketTimeoutTimestamp
			}

			msg := types.NewMsgRegisterAccount(
				clientCtx.GetFromAddress().String(),
				argPort,
				argChannel,
				timeoutTimestamp,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint64(flagTimeoutTimestamp, 0, "Timeout timestamp in nanoseconds. Default is 10 minutes.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
