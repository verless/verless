package core

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/verless/verless/builder"
	"github.com/verless/verless/config"
	"github.com/verless/verless/core/build"
	"github.com/verless/verless/parser"
	"github.com/verless/verless/plugin/atom"
	"github.com/verless/verless/plugin/tags"
	"github.com/verless/verless/writer"
)

var (
	// ErrCannotOverwrite states that verless isn't allowed to
	// delete or overwrite the output directory.
	ErrCannotOverwrite = errors.New(`Cannot overwrite the output directory.
Consider using the --overwrite flag or enabled build.overwrite in the configuration file.`)
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
		w       = writer.New(path, out)
		plugins = make([]build.Plugin, 0)
	)

	if cfg.Version == "" {
		return []error{
			errors.New("the configuration has to include a version"),
		}
	}

	if !canOverwrite(out, &options, &cfg) {
		return []error{ErrCannotOverwrite}
	}

	if cfg.HasPlugin(atom.Key) {
		atomPlugin := atom.New(&cfg.Site.Meta, out)
		plugins = append(plugins, atomPlugin)
	}

	if cfg.HasPlugin(tags.Key) {
		tagsPlugin := tags.New()
		plugins = append(plugins, tagsPlugin)
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
	if options.OutputDir != "" {
		return options.OutputDir
	}

	return filepath.Join(path, config.OutputDir)
}

// canOverwrite determines of the output directory can be removed
// or overwritten safely. This is true if
//	- the user specified the --overwrite flag,
//	- the user opted in to overwriting in verless.yml,
//	- or the output directory doesn't exist yet.
func canOverwrite(outputDir string, options *BuildOptions, cfg *config.Config) bool {
	if options.Overwrite || cfg.Build.Overwrite {
		return true
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return true
	}
	return false
}
