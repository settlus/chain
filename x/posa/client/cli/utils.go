package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/settlus/chain/x/posa/types"
)

// ParseCreateValidatorData reads and parses a ParseCreateValidatorData from a given file.
func ParseCreateValidatorData(fileName string) (*types.ValidatorInfo, error) {
	validatorInfo := types.ValidatorInfo{}

	data, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &validatorInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal validator info: %w", err)
	}

	return &validatorInfo, nil
}
