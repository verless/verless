package cli

import (
	"github.com/spf13/cobra"
	"github.com/verless/verless/core"
)

func newCreateCmd() *cobra.Command {
	createCmd := cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	createCmd.AddCommand(newCreateProjectCmd())

	return &createCmd
}

func newCreateProjectCmd() *cobra.Command {
	var (
		options core.CreateProjectOptions
	)

	createProjectCmd := cobra.Command{
		Use:  "project <path>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return core.RunCreateProject(path, options)
		},
	}

	createProjectCmd.Flags().BoolVarP(&options.Force, "force", "f",
		false, `overwrite the directory if it already exists`)

	return &createProjectCmd
}
