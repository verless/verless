package core

import (
	"path/filepath"

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
	// config.OutputDir will be used.
	OutputDir string
	// Overwrite specifies that the output folder can be overwritten.
	Overwrite bool
}

// RunBuild triggers a build using the provided options and user
// configuration.
//
// See doc.go for more information on the core architecture.
func RunBuild(path string, options BuildOptions, cfg config.Config) []error {
	var (
		out     = finalOutputDir(path, &options)
		p       = parser.NewMarkdown()
		b       = builder.New(&cfg)
		w, err  = writer.New(path, out, options.Overwrite)
		plugins = make([]build.Plugin, 0)
	)

	if err != nil {
		return []error{err}
	}

	if cfg.HasPlugin(atom.Key) {
		atomPlugin := atom.New(&cfg.Site.Meta, out)
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

// finalOutputDir determines the final output path.
func finalOutputDir(path string, options *BuildOptions) string {
	var outputPath string

	if options.OutputDir != "" {
		outputPath = options.OutputDir
	} else {
		outputPath = filepath.Join(path, config.OutputDir)
	}

	return outputPath
}
