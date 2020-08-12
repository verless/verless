package core

import (
	"github.com/verless/verless/builder"
	"github.com/verless/verless/config"
	"github.com/verless/verless/core/build"
	"github.com/verless/verless/parser"
	"github.com/verless/verless/plugin/atom"
	"github.com/verless/verless/writer"
)

// BuildOptions represents options for running a verless build.
type BuildOptions struct {
	// OutputDir sets the output directory. If this field is empty,
	// it defaults to config.OutputDir.
	OutputDir string
	// RenderRSS renders an Atom RSS feed.
	RenderRSS bool
}

// RunBuild triggers a build using the provided options and user
// configuration.
//
// It is responsible for initializing all build dependencies like
// the parser, builder and writer, which are then passed to the
// core build function. Also, all build plugins are initialized.
//
// See doc.go for more information on the core architecture.
//
// As some parts are running concurrently several errors can occur at the same time.
func RunBuild(path string, options BuildOptions, cfg config.Config) []error {
	var (
		p       = parser.NewMarkdown()
		b       = builder.New()
		w, err  = writer.New(path, options.OutputDir)
		plugins = make([]build.Plugin, 0)
	)

	if err != nil {
		return []error{err}
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
