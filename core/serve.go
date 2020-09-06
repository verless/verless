package core

import (
	"net"
	"os"

	"github.com/verless/verless/config"
	"github.com/verless/verless/core/serve"
)

// ServeOptions represents options for running a verless serve command.
type ServeOptions struct {
	BuildOptions
	// Port specifies the port to run the server at.
	Port uint16

	// Build enables automatic building of the verless project before serving.
	Build bool

	// IP specifies the ip to listen on in combination with the port.
	IP net.IP
}

// RunServe
func RunServe(path string, options ServeOptions) error {
	// First check if the passed path is a verless project. (valid verless cfg)
	cfg, err := config.FromFile(path, config.Filename)
	if err != nil {
		return err
	}

	// If yes, build it if requested to do so.
	if options.Build {
		err = RunBuild(path, options.BuildOptions, cfg)
		if err != nil {
			return nil
		}
	}

	targetFiles := finalOutputDir(path, &options.BuildOptions)

	// If the target folder doesn't exist, return an error.
	if _, err := os.Stat(targetFiles); err != nil {
		return err
	}

	// Then serve it.
	return serve.Run(serve.Context{Path: targetFiles, Port: options.Port, IP: options.IP})
}
