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
	// RecompileTemplates forces a recompilation of all templates.
	RecompileTemplates bool
}

// RunBuild triggers a build using the provided options and user
// configuration.
//
// See doc.go for more information on the core architecture.
func RunBuild(fs afero.Fs, path string, options BuildOptions, cfg config.Config) error {
	var (
		outputDir = finalOutputDir(path, &options)
		p         = parser.NewMarkdown()
		b         = builder.New(&cfg)
		w         = writer.New(fs, path, outputDir, options.RecompileTemplates)
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
		RecompileTemplates: options.RecompileTemplates,
	}

	plugins := loadPlugins(&cfg, outputDir)

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
func loadPlugins(cfg *config.Config, outputDir string) map[string]func() build.Plugin {

	plugins := map[string]func() build.Plugin{
		"atom": func() build.Plugin { return atom.New(&cfg.Site.Meta, outputDir) },
		"tags": func() build.Plugin { return tags.New() },
	}

	return plugins
}
