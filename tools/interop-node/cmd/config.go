package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	cfg "github.com/settlus/chain/tools/interop-node/config"
)

const (
	flagOverwrite   = "overwrite"
	DefaultLogLevel = "info"
	DefaultDBHome   = "db"
	DefaultPort     = 8000
)

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Commands to configure the oracle feeder",
	}

	cmd.AddCommand(initCmd())

	return cmd
}

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "initialize configuration file and home directory if one doesn't already exist",
		Long: `initialize configuration file.
for oracle feeder.
		`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cmdFlags := cmd.Flags()

			overwrite, _ := cmdFlags.GetBool(flagOverwrite)

			if _, err := os.Stat(config.ConfigFile); !os.IsNotExist(err) && !overwrite {
				return fmt.Errorf("%s already exists. Provide the -o flag to overwrite the existing config",
					config.ConfigFile)
			}

			home := config.HomeDir
			c := cfg.Config{
				Settlus: cfg.SettlusConfig{
					ChainId:  "settlus_5371-1",
					RpcUrl:   "http://localhost:26657",
					GrpcUrl:  "http://localhost:9090",
					Insecure: true,
					GasLimit: 200000,
					Fees: cfg.Fee{
						Denom:  "asetl",
						Amount: "210000000000000",
					},
				},
				Feeder: cfg.FeederConfig{
					Topics:           "block",
					SignerMode:       cfg.Local,
					Key:              "",
					ValidatorAddress: "settlusvaloper1x0foobar",
				},
				Chains: []cfg.ChainConfig{
					{
						ChainID:   "1",
						ChainName: "Ethereum",
						ChainType: "ethereum",
						RpcUrl:    "http://localhost:8545",
					},
				},
				LogLevel: DefaultLogLevel,
				Port:     DefaultPort,
				DBHome:   path.Join(home, DefaultDBHome),
			}

			config.Config = c
			if err = os.MkdirAll(config.HomeDir, 0700); err != nil {
				return err
			}
			if err = config.WriteConfigFile(); err != nil {
				return err
			}

			fmt.Printf("Successfully initialized configuration: %s\n", config.ConfigFile)

			return nil
		},
	}

	f := cmd.Flags()
	f.BoolP(flagOverwrite, "o", false, "overwrite the existing config file")

	return cmd
}
