package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
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
	// OverwriteTemplates forces a recompilation of all templates.
	OverwriteTemplates bool
}

// RunBuild triggers a build using the provided options and user
// configuration.
//
// See doc.go for more information on the core architecture.
func RunBuild(fs afero.Fs, path string, options BuildOptions, cfg config.Config) error {
	outputDir := getOutputDir(path, &options)

	writerCtx := writer.Context{
		Fs:                 fs,
		Path:               path,
		OutputDir:          outputDir,
		Theme:              cfg.Theme,
		OverwriteTemplates: options.OverwriteTemplates,
	}

	var (
		p = parser.NewMarkdown()
		b = builder.New(&cfg)
		w = writer.New(writerCtx)
	)

	if cfg.Version == "" {
		return errors.New("the configuration has to include the version key")
	}

	if !canOverwrite(fs, outputDir, &options, &cfg) {
		return ErrCannotOverwrite
	}

	ctx := build.Context{
		Path:               path,
		Parser:             p,
		Builder:            b,
		Writer:             w,
		Plugins:            make([]build.Plugin, 0),
		Types:              cfg.Types,
		RecompileTemplates: options.OverwriteTemplates,
	}

	plugins := loadPlugins(&cfg, fs, outputDir)

	for _, key := range cfg.Plugins {
		if _, exists := plugins[key]; !exists {
			return fmt.Errorf("plugin %s not found", key)
		}
		ctx.Plugins = append(ctx.Plugins, plugins[key]())
	}

	errs := build.Run(ctx)

	if len(errs) == 1 {
		return errs[0]
	} else if len(errs) > 1 {
		return errors.Errorf("several errors occurred while building: %v", errs)
	}

	return nil
}

// getOutputDir determines the final output path.
func getOutputDir(path string, options *BuildOptions) string {
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
func canOverwrite(fs afero.Fs, outputDir string, options *BuildOptions, cfg *config.Config) bool {
	if options.Overwrite || cfg.Build.Overwrite {
		return true
	}
	if _, err := fs.Stat(outputDir); os.IsNotExist(err) {
		return true
	}
	return false
}

// loadPlugins returns a map of all available plugins. Each entry
// is a function that returns a fully initialized plugin instance.
func loadPlugins(cfg *config.Config, fs afero.Fs, outputDir string) map[string]func() build.Plugin {

	plugins := map[string]func() build.Plugin{
		"atom": func() build.Plugin { return atom.New(&cfg.Site.Meta, fs, outputDir) },
		"tags": func() build.Plugin { return tags.New() },
	}

	return plugins
}
