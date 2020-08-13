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
		path    string
	)

	buildCmd := cobra.Command{
		Use: "build",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.FromFile(path, config.FileName)
			if err != nil {
				return err
			}

			errs := core.RunBuild(path, options, cfg)

			if len(errs) == 1 {
				return errs[0]
			} else if len(errs) > 1 {
				return errors.Errorf("several errors occurred while building: %v", errs)
			}

			return nil
		},
	}

	buildCmd.Flags().StringVarP(&options.OutputDir, "output", "o",
		config.OutputDir, `Specify an output directory.`)
	buildCmd.Flags().StringVarP(&path, "path", "p",
		".", `Specify a build path other than the current directory.`)
	buildCmd.Flags().BoolVar(&options.RenderRSS, "render-rss",
		true, `Render an Atom RSS feed as atom.xml.`)

	return &buildCmd
}
