package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/verless/verless/config"
	"github.com/verless/verless/core"
)

// newBuildCmd creates the `verless build` command.
func newBuildCmd() *cobra.Command {
	var (
		options core.BuildOptions
	)

	buildCmd := cobra.Command{
		Use: "build PROJECT",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.FromFile(args[0], config.Filename)
			if err != nil {
				return err
			}

			errs := core.RunBuild(args[0], options, cfg)

			if len(errs) == 1 {
				return errs[0]
			} else if len(errs) > 1 {
				return errors.Errorf("several errors occurred while building: %v", errs)
			}

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	buildCmd.Flags().StringVarP(&options.OutputDir, "output", "o",
		"", `Specify an output directory.`)

	// Overwrite should not have a shorthand to avoid accidental usage.
	buildCmd.Flags().BoolVar(&options.Overwrite, "overwrite", false, `Allows overwriting an existing output directory.`)

	return &buildCmd
}
