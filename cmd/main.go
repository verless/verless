package main

import (
	"log"

	"github.com/verless/verless/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
