package cli

import (
	"github.com/spf13/cobra"
	"github.com/verless/verless/config"
	"github.com/verless/verless/core"
)

func newBuildCmd() *cobra.Command {
	var (
		options core.BuildOptions
		path    string
	)

	buildCmd := cobra.Command{
		Use: "build",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.FromFile(path, config.ConfigName)
			if err != nil {
				return err
			}
			return core.RunBuild(path, options, cfg)
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
