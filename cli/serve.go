package cli

import (
	"net"

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

	serveCmd.Flags().BoolVarP(&options.Watch, "watch", "w",
		false, `build the project whenever something in the project changes, allows using all flags which are valid for verless build`)

	serveCmd.Flags().IPVarP(&options.IP, "ip", "i",
		net.IP{0, 0, 0, 0}, `specify the IP to listen on, it has to be a valid IPv4 or IPv6`)

	addBuildOptions(&serveCmd, &options.BuildOptions, false)

	return &serveCmd
}
