package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	cfg "github.com/settlus/chain/tools/interop-node/config"
)

var config cfg.RuntimeConfig

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interop-node",
		Short: "interop-node - Settlus interoperability helper node",
	}

	cmd.AddCommand(configCmd())
	cmd.AddCommand(startCmd())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var home string
	if config.HomeDir == "" {
		userHome, err := homedir.Dir()
		handleInitError(err)
		home = filepath.Join(userHome, ".interop")
	} else {
		home = config.HomeDir
	}

	config = cfg.RuntimeConfig{
		HomeDir:    home,
		ConfigFile: filepath.Join(home, "config.yaml"),
	}
	viper.SetConfigFile(config.ConfigFile)
	viper.SetEnvPrefix("interop")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("no config exists at default location", err)
		return
	}
	handleInitError(viper.Unmarshal(&config.Config))
	bz, err := os.ReadFile(viper.ConfigFileUsed())
	handleInitError(err)
	handleInitError(yaml.Unmarshal(bz, &config.Config))

	handleInitError(config.Config.Validate())
}

func handleInitError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
