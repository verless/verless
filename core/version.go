package core

import (
	"fmt"

	"github.com/verless/verless/config"
)

// VersionOptions represents options for the version command.
type VersionOptions struct {
	// Quiet only prints the plain version number.
	Quiet bool
}

// RunVersion prints verless version information.
func RunVersion(options VersionOptions) error {
	var format string

	if options.Quiet {
		format = "%s"
	} else {
		format = "verless version %s"
	}

	fmt.Printf(format, config.Version)
	return nil
}
