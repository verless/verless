package core

import (
	"github.com/verless/verless/builder"
	"github.com/verless/verless/config"
	"github.com/verless/verless/core/build"
	"github.com/verless/verless/parser"
	"github.com/verless/verless/plugin/atom"
	"github.com/verless/verless/writer"
)

type BuildOptions struct {
	OutputDir string
	RenderRSS bool
}

func RunBuild(path string, options BuildOptions, cfg config.Config) error {
	var (
		p       = parser.NewMarkdown()
		b       = builder.New()
		w, err  = writer.New(path, options.OutputDir)
		plugins = make([]build.Plugin, 0)
	)

	if err != nil {
		return err
	}

	if options.RenderRSS {
		atomPlugin := atom.New(&cfg.Site.Meta, options.OutputDir)
		plugins = append(plugins, atomPlugin)
	}

	ctx := build.Context{
		Path:    path,
		Parser:  p,
		Builder: b,
		Writer:  w,
		Plugins: plugins,
	}

	return build.Run(ctx)
}
