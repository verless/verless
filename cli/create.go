package cli

import (
	"github.com/spf13/cobra"
	"github.com/verless/verless/core"
)

// newCreateCmd creates the `verless create` command.
func newCreateCmd() *cobra.Command {
	createCmd := cobra.Command{
		Use:   "create",
		Short: `Create a new verless object`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	createCmd.AddCommand(newCreateProjectCmd())
	createCmd.AddCommand(newCreateThemeCmd())
	createCmd.AddCommand(newCreateFile())

	return &createCmd
}

// newCreateFile creates the `verless create file` command
func newCreateFile() *cobra.Command {
	var (
		options core.CreateFileOptions
	)
	createFileCmd := cobra.Command{
		Use:   "file NAME",
		Short: `Create a new content file`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return core.CreateFile(path, options)
		},
	}

	createFileCmd.Flags().StringVarP(&options.Project, "project", "p", ".", `project path to create file in.`)

	return &createFileCmd
}

// newCreateProjectCmd creates the `verless create project` command.
func newCreateProjectCmd() *cobra.Command {
	var (
		options core.CreateProjectOptions
	)

	createProjectCmd := cobra.Command{
		Use:   "project NAME",
		Short: `Create a new verless project`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return core.CreateProject(path, options)
		},
	}

	createProjectCmd.Flags().BoolVar(&options.Overwrite, "overwrite",
		false, `overwrite the directory if it already exists`)

	return &createProjectCmd
}

// newCreateThemeCmd creates the `verless create theme` command.
func newCreateThemeCmd() *cobra.Command {
	createThemeCmd := cobra.Command{
		Use:   "theme PROJECT NAME",
		Short: `Create a new verless theme`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			name := args[1]

			return core.CreateTheme(path, name)
		},
	}

	return &createThemeCmd
}
