package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"gopkg.in/yaml.v2"
)

const (
	AwsKms = "aws-kms"
	Local  = "local"
)

type RuntimeConfig struct {
	HomeDir    string
	ConfigFile string
	Config     Config
}

func (c RuntimeConfig) WriteConfigFile() error {
	return os.WriteFile(c.ConfigFile, c.Config.MustMarshalYaml(), 0600)
}

type Config struct {
	Settlus  SettlusConfig `yaml:"settlus"`
	Feeder   FeederConfig  `yaml:"feeder"`
	Chains   []ChainConfig `yaml:"chains"`
	DBHome   string        `yaml:"db_home"`
	LogLevel string        `yaml:"log_level"`
}

func (c *Config) Validate() error {
	for _, chain := range c.Chains {
		if err := chain.Validate(); err != nil {
			return err
		}
	}

	if err := c.Settlus.Validate(); err != nil {
		return err
	}

	if err := c.Feeder.Validate(); err != nil {
		return err
	}

	if len(c.DBHome) == 0 {
		return fmt.Errorf("db_home must not be empty")
	}

	if len(c.LogLevel) == 0 {
		return fmt.Errorf("log_level must not be empty")
	}

	return nil
}

func (c *Config) MustMarshalYaml() []byte {
	out, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return out
}

type SettlusConfig struct {
	ChainId  string `yaml:"chain_id"`
	RpcUrl   string `yaml:"rpc_url"`
	GrpcUrl  string `yaml:"grpc_url"`
	Insecure bool   `yaml:"insecure"`
	GasLimit uint64 `yaml:"gas_limit"`
	Fees     Fee    `yaml:"fees"`
}

type Fee struct {
	Denom  string `yaml:"denom"`
	Amount string `yaml:"amount"`
}

func (sc *SettlusConfig) Validate() error {
	if len(sc.ChainId) == 0 {
		return fmt.Errorf("settlus chain_id must not be empty")
	}

	if _, err := url.Parse(sc.RpcUrl); err != nil {
		return fmt.Errorf("invalid settlus rpc_url: %w", err)
	}

	if _, err := url.Parse(sc.GrpcUrl); err != nil {
		return fmt.Errorf("invalid settlus grpc_url: %w", err)
	}

	if sc.GasLimit <= 0 {
		return fmt.Errorf("settlus gas_limit must be larger than 0: %d", sc.GasLimit)
	}

	amount, ok := sdk.NewIntFromString(sc.Fees.Amount)
	if !ok {
		return fmt.Errorf("invalid amount %s", sc.Fees.Amount)
	}

	fees := sdk.NewCoin(sc.Fees.Denom, amount)

	if fees.IsNegative() || fees.IsZero() {
		return fmt.Errorf("invalid fees %s", fees.String())
	}

	return nil
}

type FeederConfig struct {
	Topics           string `yaml:"topics"`
	Address          string `yaml:"address"`
	SignerMode       string `yaml:"signer_mode"`
	Key              string `yaml:"key"`
	ValidatorAddress string `yaml:"validator_address"`
}

func (fc *FeederConfig) Validate() error {
	if fc.Topics == "" {
		return fmt.Errorf("at least one topic must be provided")
	}

	if fc.ValidatorAddress == "" {
		return fmt.Errorf("validator_address must not be empty")
	}

	if !strings.HasPrefix(fc.ValidatorAddress, "setlvaloper") {
		return fmt.Errorf("validator_address must start with setlvaloper: %s", fc.ValidatorAddress)
	}

	if fc.SignerMode != AwsKms && fc.SignerMode != Local {
		return fmt.Errorf("invalid signer_mode, must be one of: %s, %s", AwsKms, Local)
	}

	if fc.Address == "" {
		return fmt.Errorf("address must not be empty")
	}

	return nil
}

type ChainConfig struct {
	ChainID   string `yaml:"chain_id"`
	ChainName string `yaml:"chain_name"`
	ChainType string `yaml:"chain_type"`
	RpcUrl    string `yaml:"rpc_url"`
}

func (c *ChainConfig) Validate() error {
	if len(c.ChainID) == 0 {
		return fmt.Errorf("chain_id must not be empty")
	}

	if len(c.ChainName) == 0 {
		return fmt.Errorf("chain_name must not be empty")
	}

	_, err := url.Parse(c.RpcUrl)
	return err
}
