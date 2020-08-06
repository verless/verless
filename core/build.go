package core

import (
	"github.com/verless/verless/builder"
	"github.com/verless/verless/config"
	"github.com/verless/verless/core/build"
	"github.com/verless/verless/parser"
	"github.com/verless/verless/writer"
)

type BuildOptions struct {
	OutputDir string
	RenderRSS bool
}

func RunBuild(path string, options BuildOptions, cfg config.Config) error {
	var (
		_parser  = parser.NewMarkdown()
		_builder = builder.New()
		plugins  = make([]build.Plugin, 0)
		_writer  = writer.New()
	)

	ctx := build.Context{
		Path:    path,
		Parser:  _parser,
		Builder: _builder,
		Plugins: plugins,
		Writer:  _writer,
	}

	return build.Run(ctx)
}
