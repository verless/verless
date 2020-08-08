package main

import (
	"github.com/verless/verless/config"
	"github.com/verless/verless/core"
)

func main() {
	// this is temporary code just to test running the app

	options := core.BuildOptions{
		OutputDir: "out",
		RenderRSS: false,
	}

	cfg, err := config.FromFile(".", "config.json")
	if err != nil {
		panic(err)
	}
	err = core.RunBuild("./in", options, cfg)
	if err != nil {
		panic(err)
	}
}
