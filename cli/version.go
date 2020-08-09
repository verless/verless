package cli

import (
	"github.com/spf13/cobra"
	"github.com/verless/verless/core"
)

func newVersionCmd() *cobra.Command {
	var options core.VersionOptions

	versionCmd := cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return core.RunVersion(options)
		},
	}

	versionCmd.Flags().BoolVarP(&options.Quiet, "quiet", "q",
		false, `Only print the version number.`)

	return &versionCmd
}
