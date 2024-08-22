package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/settlus/chain/tools/interop-node/version"
)

const (
	flagLong   = "long"
	flagOutput = "output"
)

// VersionCommand returns a command that prints the application binary version information.
func VersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the application binary version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			verInfo := version.NewInfo()

			long, _ := cmd.Flags().GetBool(flagLong)
			output, _ := cmd.Flags().GetString(flagOutput)

			if !long {
				fmt.Fprintln(cmd.OutOrStdout(), verInfo.Version)
				return nil
			}

			verInfo.ExtraInfo = extraInfoFromContext(cmd)

			var (
				bz  []byte
				err error
			)

			switch strings.ToLower(output) {
			case "json":
				bz, err = json.MarshalIndent(createSimplifiedVersion(verInfo), "", "  ")
			default:
				bz, err = yaml.Marshal(createSimplifiedVersion(verInfo))
			}

			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(bz))
			return nil
		},
	}

	cmd.Flags().Bool(flagLong, false, "Print long version information")
	cmd.Flags().StringP(flagOutput, "o", "text", "Output format (text|json)")

	return cmd
}

func createSimplifiedVersion(info version.Info) map[string]interface{} {
	return map[string]interface{}{
		"name":       info.Name,
		"version":    info.Version,
		"go_version": info.GoVersion,
		"build_tags": info.BuildTags,
		"extra_info": info.ExtraInfo,
	}
}

func extraInfoFromContext(cmd *cobra.Command) version.ExtraInfo {
	ctx := cmd.Context()
	if ctx != nil {
		extraInfo, ok := ctx.Value(version.ContextKey{}).(version.ExtraInfo)
		if ok {
			return extraInfo
		}
	}
	return nil
}
