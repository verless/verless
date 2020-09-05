package cli

import (
	"github.com/spf13/cobra"
	"github.com/verless/verless/core"
)

// newServeCmd creates the `verless serve` command.
func newServeCmd() *cobra.Command {
	var (
		options core.ServeOptions
	)

	serveCmd := cobra.Command{
		Use: "serve PROJECT",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := core.RunServe(args[0], options)

			if err != nil {
				return err
			}

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	serveCmd.Flags().Uint16VarP(&options.Port, "port", "p",
		8080, `specify the port for the web server`)

	serveCmd.Flags().BoolVarP(&options.Build, "build", "b",
		false, `build the project before serving, allows using all flags which are valid for verless build`)

	addBuildOptions(&serveCmd, &options.BuildOptions)

	return &serveCmd
}
