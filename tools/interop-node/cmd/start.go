package cmd

import (
	"context"
	"fmt"
	"net"

	"github.com/tendermint/tendermint/libs/cli/flags"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/settlus/chain/tools/interop-node/server"
	"github.com/settlus/chain/x/interop"
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

			s := grpc.NewServer()
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Config.Port))
			if err != nil {
				return err
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

			interop.RegisterInteropServer(s, interopServer)
			err = s.Serve(lis)
			if err != nil {
				return fmt.Errorf("failed to start grpc service: %w", err)
			}

			return nil
		},
	}

	return cmd
}
