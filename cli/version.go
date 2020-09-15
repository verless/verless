package cli

import (
	"github.com/spf13/cobra"
	"github.com/verless/verless/core"
)

// newVersionCmd creates the `verless version` command.
func newVersionCmd() *cobra.Command {
	var options core.VersionOptions

	versionCmd := cobra.Command{
		Use:   "version",
		Short: `Print version information`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return core.RunVersion(options)
		},
	}

	versionCmd.Flags().BoolVarP(&options.Quiet, "quiet", "q",
		false, `only print the version number`)

	return &versionCmd
}
