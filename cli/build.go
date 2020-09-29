package cli

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/verless/verless/core"
)

// newBuildCmd creates the `verless build` command.
func newBuildCmd() *cobra.Command {
	var (
		options core.BuildOptions
	)

	buildCmd := cobra.Command{
		Use:   "build PROJECT",
		Short: `Build your verless project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			targetFs := afero.NewOsFs()

			build, err := core.NewBuild(targetFs, path, options)
			if err != nil {
				return err
			}

			return build.Run()
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
