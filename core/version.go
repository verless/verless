package core

import (
	"fmt"

	"github.com/verless/verless/config"
)

var (
	// format is the default format string for printing the version.
	format = `verless version %s
Git tag: %s
Git commit: %s`
)

// VersionOptions represents options for the version command.
type VersionOptions struct {
	// Quiet only prints the plain version number.
	Quiet bool
}

// RunVersion prints verless version information.
func RunVersion(options VersionOptions) error {
	if options.Quiet {
		fmt.Printf("%s", config.GitTag)
	}

	fmt.Printf(format, config.GitTag, config.GitTag, config.GitCommit)
	return nil
}
