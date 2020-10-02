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

	return &createCmd
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
	var (
		options core.CreateThemeOptions
	)
	createThemeCmd := cobra.Command{
		Use:   "theme THEME_NAME",
		Short: `Create a new verless theme`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			return core.CreateTheme(options, name)
		},
	}

	createThemeCmd.Flags().StringVarP(&options.Project, "project", "p", ".", `project path to create new theme in.`)
	return &createThemeCmd
}
