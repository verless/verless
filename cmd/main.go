// Package main provides the verless application.
package main

import (
	"log"

	"github.com/verless/verless/cli"
)

// main runs the verless CLI.
func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
