package core

import (
	"fmt"

	"github.com/verless/verless/config"
)

type VersionOptions struct {
	Quiet bool
}

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
