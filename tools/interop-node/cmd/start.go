package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cometbft/cometbft/libs/cli/flags"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/spf13/cobra"

	"github.com/settlus/chain/tools/interop-node/server"
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "start",
		Short:        "Start interop node",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			logger := log.NewTMLogger(log.NewSyncWriter(out))
			logger, err := flags.ParseLogLevel(config.Config.LogLevel, logger, DefaultLogLevel)
			if err != nil {
				return fmt.Errorf("failed to parse log level: %w", err)
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			interopServer, err := server.NewServer(&config.Config, ctx, logger)
			if err != nil {
				cancel()
				return fmt.Errorf("failed to create interOp server: %w", err)
			}

			interopServer.Start()
			defer func() {
				cancel()
				interopServer.Close()
			}()

			cancelChan := make(chan os.Signal, 1)
			signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

			<-cancelChan

			return nil
		},
	}

	return cmd
}
