package e2e

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"time"
)

const (
	flagFrom            = "from"
	flagHome            = "home"
	flagFees            = "fees"
	flagGas             = "gas"
	flagOutput          = "output"
	flagChainID         = "chain-id"
	flagSpendLimit      = "spend-limit"
	flagGasAdjustment   = "gas-adjustment"
	flagFeeAccount      = "fee-account"
	flagBroadcastMode   = "broadcast-mode"
	flagKeyringBackend  = "keyring-backend"
	flagAllowedMessages = "allowed-messages"

	chainId = "settlus_5371-1"
)

type flagOption func(map[string]interface{})

// withKeyValue add a new flag to command
func withKeyValue(key string, value interface{}) flagOption {
	return func(o map[string]interface{}) {
		o[key] = value
	}
}

func applyOptions(chainID string, options []flagOption) map[string]interface{} {
	opts := map[string]interface{}{
		flagKeyringBackend: "test",
		flagOutput:         "json",
		flagGas:            "auto",
		flagFrom:           "bob",
		flagBroadcastMode:  "sync",
		flagGasAdjustment:  "1.5",
		flagChainID:        chainID,
		flagFees:           standardFees.String(),
	}
	for _, apply := range options {
		apply(opts)
	}
	return opts
}

func (s *IntegrationTestSuite) execBankSend(
	from,
	to,
	amt,
	fees string,
	opt ...flagOption,
) {
	opt = append(opt, withKeyValue(flagFees, fees))
	opt = append(opt, withKeyValue(flagFrom, from))
	opts := applyOptions(chainId, opt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	settlusCmd := []string{
		settlusdBinary,
		txCommand,
		banktypes.ModuleName,
		"send",
		from,
		to,
		amt,
		"-y",
	}
	for flag, value := range opts {
		settlusCmd = append(settlusCmd, fmt.Sprintf("--%s=%v", flag, value))
	}

	s.executeSettlusTxCommand(ctx, settlusCmd)
}

func (s *IntegrationTestSuite) execCreateTenant(
	from,
	denom string,
	period string,
	opt ...flagOption,
) {
	opt = append(opt, withKeyValue(flagFees, "1010000uusdc"))
	opt = append(opt, withKeyValue(flagFrom, from))
	opts := applyOptions(chainId, opt)
	opts[flagGas] = "1010000"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	settlusCmd := []string{
		settlusdBinary,
		txCommand,
		"settlement",
		"create-tenant",
		denom,
		period,
		"-y",
	}
	for flag, value := range opts {
		settlusCmd = append(settlusCmd, fmt.Sprintf("--%s=%v", flag, value))
	}

	s.executeSettlusTxCommand(ctx, settlusCmd)
}

func (s *IntegrationTestSuite) execCreateMcTenant(
	from,
	denom string,
	period string,
	opt ...flagOption,
) {
	opt = append(opt, withKeyValue(flagFees, "1010000uusdc"))
	opt = append(opt, withKeyValue(flagFrom, from))
	opts := applyOptions(chainId, opt)
	opts[flagGas] = "1010000"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	settlusCmd := []string{
		settlusdBinary,
		txCommand,
		"settlement",
		"create-tenant-mc",
		denom,
		period,
		"-y",
	}
	for flag, value := range opts {
		settlusCmd = append(settlusCmd, fmt.Sprintf("--%s=%v", flag, value))
	}

	s.executeSettlusTxCommand(ctx, settlusCmd)
}

func (s *IntegrationTestSuite) execRecord(
	from string,
	tenantId uint64,
	requestId,
	amount,
	extChainId,
	contractAddr,
	tokenIdHex string,
	opt ...flagOption,
) {
	opt = append(opt, withKeyValue(flagFees, "10000uusdc"))
	opt = append(opt, withKeyValue(flagFrom, from))
	opts := applyOptions(chainId, opt)
	opts[flagGas] = "10000"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	settlusCmd := []string{
		settlusdBinary,
		txCommand,
		"settlement",
		"record",
		strconv.Itoa(int(tenantId)),
		requestId,
		amount,
		extChainId,
		contractAddr,
		tokenIdHex,
		"",
		"-y",
	}

	for flag, value := range opts {
		settlusCmd = append(settlusCmd, fmt.Sprintf("--%s=%v", flag, value))
	}

	s.executeSettlusTxCommand(ctx, settlusCmd)
}

func (s *IntegrationTestSuite) execDepositToTreasury(
	from string,
	tenantId uint64,
	amount string,
	opt ...flagOption,
) {
	opt = append(opt, withKeyValue(flagFees, "10000uusdc"))
	opt = append(opt, withKeyValue(flagFrom, from))
	opts := applyOptions(chainId, opt)
	opts[flagGas] = "10000"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	settlusCmd := []string{
		settlusdBinary,
		txCommand,
		"settlement",
		"deposit-to-treasury",
		strconv.Itoa(int(tenantId)),
		amount,
		"-y",
	}

	for flag, value := range opts {
		settlusCmd = append(settlusCmd, fmt.Sprintf("--%s=%v", flag, value))
	}

	s.executeSettlusTxCommand(ctx, settlusCmd)
}

func (s *IntegrationTestSuite) executeSettlusTxCommand(ctx context.Context, settlusCmd []string) {
	cmd := exec.CommandContext(ctx, settlusCmd[0], settlusCmd[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		s.T().Log(">>", cmd.String())
		s.T().Log("<<", out.String())
	}
	s.Require().NoError(err)
}
