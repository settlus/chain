package types

import (
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// constants
const (
	ProposalTypeCreateValidator string = "CreateValidatorProposal" // #nosec
)

// Implements Proposal Interface
var (
	_ v1beta1.Content = &CreateValidatorProposal{}
)

func init() {
	v1beta1.RegisterProposalType(ProposalTypeCreateValidator)
	v1beta1.ModuleCdc.Amino.RegisterConcrete(&CreateValidatorProposal{}, "poa/CreateValidatorProposal", nil)
}

// NewCreateValidatorProposal returns new instance of CreateValidatorProposal
func NewCreateValidatorProposal(title, description string, info *ValidatorInfo) v1beta1.Content {
	return &CreateValidatorProposal{
		Title:       title,
		Description: description,
		Info:        info,
	}
}

// ProposalRoute returns router key for this proposal
func (*CreateValidatorProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*CreateValidatorProposal) ProposalType() string {
	return ProposalTypeCreateValidator
}

// ValidateBasic performs a stateless check of the proposal fields
func (ttcp *CreateValidatorProposal) ValidateBasic() error {
	// check if the token is a hex address, if not, check if it is a valid SDK
	// // denom
	return v1beta1.ValidateAbstract(ttcp)
}
