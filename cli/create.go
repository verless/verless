package cli

import (
	"github.com/spf13/cobra"
	"github.com/verless/verless/core"
)

// newCreateCmd creates the `verless create` command.
func newCreateCmd() *cobra.Command {
	createCmd := cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	createCmd.AddCommand(newCreateProjectCmd())
	createCmd.AddCommand(newCreateExampleCmd())

	return &createCmd
}

// newCreateProjectCmd creates the `verless create project` command.
func newCreateProjectCmd() *cobra.Command {
	var (
		options core.CreateProjectOptions
	)

	createProjectCmd := cobra.Command{
		Use:  "project NAME",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return core.RunCreateProject(path, options)
		},
	}

	createProjectCmd.Flags().BoolVar(&options.Overwrite, "overwrite",
		false, `overwrite the directory if it already exists`)

	return &createProjectCmd
}

func newCreateExampleCmd() *cobra.Command {
	var (
		options core.CreateExampleOptions
	)

	createExampleCmd := cobra.Command{
		Use: "example",
		RunE: func(cmd *cobra.Command, args []string) error {
			return core.RunCreateExample(options)
		},
	}

	createExampleCmd.Flags().BoolVar(&options.Overwrite, "overwrite",
		false, `overwrite the directory if it already exists`)

	return &createExampleCmd
}
