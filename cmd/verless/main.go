// Package main provides the verless application.
package main

import (
	"github.com/verless/verless/cli"
	"github.com/verless/verless/out"
	"github.com/verless/verless/out/style"
)

// main runs the verless CLI.
func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		out.T(style.X, err.Error())
	}
}
