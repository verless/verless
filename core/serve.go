package core

import "github.com/verless/verless/core/serve"

// ServeOptions represents options for running a verless serve command.
type ServeOptions struct {
	// Port specifies the port to run the server at.
	Port uint16
}

// RunServe
func RunServe(path string, options ServeOptions) []error {
	return serve.Run(serve.Context{Path: path})
}
