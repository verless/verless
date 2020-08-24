package cli

import (
	"github.com/spf13/cobra"
	"github.com/verless/verless/config"
)

// NewRootCmd creates the `verless` command and its sub-commands.
func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:     "verless",
		Short:   `A simple and lightweight Static Site Generator.`,
		Version: config.GitTag,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(newBuildCmd())
	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newVersionCmd())

	return &rootCmd
}
