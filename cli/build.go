package cli

import (
	"github.com/spf13/afero"
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

			err = core.RunBuild(afero.NewOsFs(), args[0], options, cfg)
			return err
		},
		Args: cobra.ExactArgs(1),
	}

	addBuildOptions(&buildCmd, &options, true)

	return &buildCmd
}

func addBuildOptions(buildCmd *cobra.Command, options *core.BuildOptions, addOverwrite bool) {
	buildCmd.Flags().StringVarP(&options.OutputDir, "output", "o",
		"", `specify an output directory`)

	if addOverwrite {
		// Overwrite should not have a shorthand to avoid accidental usage.
		buildCmd.Flags().BoolVar(&options.Overwrite, "overwrite",
			false, `allows overwriting an existing output directory`)
	}
}
