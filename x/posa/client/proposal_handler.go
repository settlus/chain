package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/settlus/chain/x/posa/client/cli"
)

// CreateValidatorProposalHandler is the new create validator proposal handler.
var (
	CreateValidatorProposalHandler = govclient.NewProposalHandler(cli.NewCreateValidatorProposalCmd)
)
