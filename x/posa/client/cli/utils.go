package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/settlus/chain/x/posa/types"
)

// ParseCreateValidatorData reads and parses a ParseCreateValidatorData from a file.
func ParseCreateValidatorData(metadataFile string) (*types.ValidatorInfo, error) {
	validatorInfo := types.ValidatorInfo{}

	contents, err := os.ReadFile(filepath.Clean(metadataFile))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(contents, &validatorInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proposal metadata: %w", err)
	}

	return &validatorInfo, nil
}
