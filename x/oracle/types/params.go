package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var (
	KeyVotePeriod                 = []byte("VotePeriod")
	KeyVoteThreshold              = []byte("VoteThreshold")
	KeyToleratedErrorBand         = []byte("ToleratedErrorBand")
	KeyWhitelist                  = []byte("Whitelist")
	KeySlashFraction              = []byte("SlashFraction")
	KeySlashWindow                = []byte("SlashWindow")
	KeyMaxMissCountPerSlashWindow = []byte("MaxMissCountPerSlashWindow")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

var (
	DefaultVotePeriod                 = uint64(10)
	DefaultVoteThreshold              = sdk.NewDecWithPrec(50, 2) // 50%
	DefaultSlashFraction              = sdk.NewDecWithPrec(1, 2)  // 1%
	DefaultSlashWindow                = uint64(100000)
	DefaultMaxMissCountPerSlashWindow = uint64(60)
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		VotePeriod:                 DefaultVotePeriod,
		VoteThreshold:              DefaultVoteThreshold,
		SlashFraction:              DefaultSlashFraction,
		SlashWindow:                DefaultSlashWindow,
		MaxMissCountPerSlashWindow: DefaultMaxMissCountPerSlashWindow,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyVotePeriod, &p.VotePeriod, validateVotePeriod),
		paramtypes.NewParamSetPair(KeyVoteThreshold, &p.VoteThreshold, validateVoteThreshold),
		paramtypes.NewParamSetPair(KeySlashFraction, &p.SlashFraction, validateSlashFraction),
		paramtypes.NewParamSetPair(KeySlashWindow, &p.SlashWindow, validateSlashWindow),
		paramtypes.NewParamSetPair(KeyMaxMissCountPerSlashWindow, &p.MaxMissCountPerSlashWindow, validateMaxMissCountPerSlashWindow),
	}
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.VotePeriod == 0 {
		return ErrVotePeriodIsZero
	}

	if p.VoteThreshold.LT(sdk.NewDecWithPrec(50, 2)) {
		return errorsmod.Wrapf(ErrInvalidParams, "vote threshold must be bigger than 50%%: %s", p.VoteThreshold)
	}

	if p.VoteThreshold.GT(sdk.NewDec(1)) {
		return errorsmod.Wrapf(ErrInvalidParams, "vote threshold %s is greater than 1", p.VoteThreshold)
	}

	if p.SlashFraction.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidParams, "slash fraction must be positive: %s", p.SlashFraction)
	}

	if p.SlashFraction.GT(sdk.NewDec(1)) {
		return errorsmod.Wrapf(ErrInvalidParams, "slash fraction %s is greater than 1", p.SlashFraction)
	}

	if p.SlashWindow <= 0 {
		return errorsmod.Wrapf(ErrInvalidParams, "slash slashWindow %d is less than or equal to 0", p.SlashWindow)
	}

	if p.VotePeriod > p.SlashWindow {
		return errorsmod.Wrapf(ErrInvalidParams, "vote period %d is greater than slash slashWindow %d", p.VotePeriod, p.SlashWindow)
	}

	if p.SlashWindow%p.VotePeriod != 0 {
		return errorsmod.Wrapf(ErrInvalidParams, "slash slashWindow %d is not divisible by vote period %d", p.SlashWindow, p.VotePeriod)
	}

	if p.SlashWindow < p.VotePeriod {
		return errorsmod.Wrapf(ErrInvalidParams, "slash slashWindow %d is less than vote period %d", p.SlashWindow, p.VotePeriod)
	}

	if p.MaxMissCountPerSlashWindow <= 0 {
		return errorsmod.Wrapf(ErrInvalidParams, "max miss count per slash slashWindow %d is less than or equal to 0", p.MaxMissCountPerSlashWindow)
	}

	if p.MaxMissCountPerSlashWindow >= p.SlashWindow {
		return errorsmod.Wrapf(ErrInvalidParams, "max miss count per slash slashWindow %d is greater than slash slashWindow %d", p.MaxMissCountPerSlashWindow, p.SlashWindow)
	}

	return nil
}

func validateVotePeriod(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidParams, "invalid parameter type: %T", i)
	}

	if v == 0 {
		return errorsmod.Wrapf(ErrInvalidParams, "vote period must be positive: %d", v)
	}

	return nil
}

func validateVoteThreshold(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidParams, "invalid parameter type: %T", i)
	}

	if v.LT(sdk.NewDecWithPrec(50, 2)) {
		return errorsmod.Wrapf(ErrInvalidParams, "vote threshold must be larger than 50%%: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return errorsmod.Wrapf(ErrInvalidParams, "vote threshold too large: %s", v)
	}

	return nil
}

func validateSlashFraction(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidParams, "invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidParams, "slash fraction must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return errorsmod.Wrapf(ErrInvalidParams, "slash fraction is too large: %s", v)
	}

	return nil
}

func validateSlashWindow(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidParams, "invalid parameter type: %T", i)
	}

	if v == 0 {
		return errorsmod.Wrapf(ErrInvalidParams, "slash window must be positive: %d", v)
	}

	return nil
}

func validateMaxMissCountPerSlashWindow(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidParams, "invalid parameter type: %T", i)
	}

	if v == 0 {
		return errorsmod.Wrapf(ErrInvalidParams, "max miss count per slash window must be positive: %d", v)
	}

	return nil
}

func CalculateRoundId(blockHeight int64, votePeriod uint64) uint64 {
	return CalculateRoundStartHeight(blockHeight, votePeriod)
}

func CalculateRoundStartHeight(blockHeight int64, votePeriod uint64) uint64 {
	uBlockHeight := uint64(blockHeight)
	return uBlockHeight - uBlockHeight%(votePeriod*2)
}

func CalculateVotePeriod(blockHeight int64, votePeriod uint64) (int64, int64) {
	iVotePeriod := int64(votePeriod)
	prevoteEnd := blockHeight - blockHeight%(iVotePeriod*2) + iVotePeriod - 1
	voteEnd := blockHeight - blockHeight%(iVotePeriod*2) + iVotePeriod*2 - 1

	return prevoteEnd, voteEnd
}
