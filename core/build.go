package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/verless/verless/builder"
	"github.com/verless/verless/config"
	"github.com/verless/verless/core/build"
	"github.com/verless/verless/fs"
	"github.com/verless/verless/model"
	"github.com/verless/verless/parser"
	"github.com/verless/verless/plugin/atom"
	"github.com/verless/verless/plugin/tags"
	"github.com/verless/verless/writer"
)

const (
	// parallelism specifies the number of parallel workers.
	parallelism int = 4
)

var (
	// ErrCannotOverwrite states that verless isn't allowed to delete or
	// overwrite the output directory.
	ErrCannotOverwrite = errors.New(`Cannot overwrite the output directory.
Consider using the --overwrite flag or enabled build.overwrite in the configuration file.`)
)

// Parser represents a parser that processes Markdown files and converts
// them into a model instance.
type Parser interface {
	// ParsePage must be safe for concurrent usage.
	ParsePage(src []byte) (model.Page, error)
}

// Builder represents a model builder that maintains a Site instance and
// registers all parsed pages in that instance.
type Builder interface {
	// RegisterPage must be safe for concurrent usage.
	RegisterPage(page model.Page) error
	Dispatch() (model.Site, error)
}

// Writer represents a model writer that renders the site model as HTML
// using the corresponding templates.
type Writer interface {
	Write(site model.Site) error
}

// Plugin represents a built-in verless plugin.
type Plugin interface {
	// ProcessPage will be invoked after parsing the page. Must be safe
	// for concurrent usage.
	ProcessPage(page *model.Page) error
	// PreWrite will be invoked before writing the site.
	PreWrite(site *model.Site) error
	// PostWrite will be invoked after writing the site.
	PostWrite() error
}

// BuildOptions represents options for running a verless build.
type BuildOptions struct {
	// OutputDir sets the output directory. If this field is empty, config.OutputDir
	// will be used.
	OutputDir string
	// Overwrite specifies that the output folder can be overwritten.
	Overwrite bool
	// RecompileTemplates forces a recompilation of all templates.
	RecompileTemplates bool
}

// Build provides methods for building a static site.
type Build struct {
	Path    string
	Parser  Parser
	Builder Builder
	Writer  Writer
	Plugins []Plugin
	Types   map[string]*model.Type
	Options BuildOptions
}

// Run executes the build using the provided build context.
//
// The current build implementation runs the following steps to build
// the static site:
//	1. Read all Markdown files and send them through a channel.
//	2. Spawn n workers reading from the channel, where n = `parallelism`.
//	3. Build the pages concurrently:
//		3.1. Read the file as a []byte
//		3.2. Parse the file and convert it to a model.Page.
//		3.3. Register the page in the builder's site model.
//		3.4. Let each plugin process the page.
//	4. Get the finished site model with all pages from the builder.
//	5. Render that site model as HTML.
//	6. Let each plugin finish its work, e.g. by writing a file.
//
// If any error occurs it is returned as a slice of errors as there can be
// several errors from the concurrent goroutines at the same time. If any
// error occurs all goroutines are stopped as soon as possible.
func (b *Build) Run(targetFs afero.Fs, path string, cfg config.Config) []error {
	outputDir := getOutputDir(b.Path, &b.Options)

	if !canOverwrite(targetFs, outputDir, &b.Options, &cfg) {
		return []error{ErrCannotOverwrite}
	}

	if cfg.Version == "" {
		return []error{errors.New("the configuration has to include the version key")}
	}

	writerCtx := writer.Context{
		Fs:                 targetFs,
		Path:               b.Path,
		OutputDir:          outputDir,
		Theme:              cfg.Theme,
		RecompileTemplates: b.Options.RecompileTemplates,
	}

	b.Path = path
	b.Parser = parser.NewMarkdown()
	b.Builder = builder.New(&cfg)
	b.Writer = writer.New(writerCtx)
	b.Types = cfg.Types

	plugins := loadPlugins(&cfg, targetFs, outputDir)

	for _, key := range cfg.Plugins {
		if _, exists := plugins[key]; !exists {
			return []error{fmt.Errorf("plugin %s not found", key)}
		}
		b.Plugins = append(b.Plugins, plugins[key]())
	}

	var (
		files      = make(chan string)
		errorChan  = make(chan error)
		retErrors  = make([]error, 0)
		contentDir = filepath.Join(b.Path, config.ContentDir)
	)

	go func() {
		if err := fs.StreamFiles(contentDir, files, fs.MarkdownOnly, fs.NoUnderscores); err != nil {
			errorChan <- err
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(parallelism)

	for i := 0; i < parallelism; i++ {
		go func() {
			// Process the files received via the files channel.
			for file := range files {
				if err := b.processFile(contentDir, file); err != nil {
					errorChan <- err
				}
			}
			wg.Done()
		}()
	}

	// Observe the WaitGroup and close the error channel when all workers
	// have finished.
	go func() {
		for {
			wg.Wait()
			close(errorChan)
			break
		}
	}()

	// Collect all errors. This automatically stops when the workers have
	// finished, as the channel gets closed.
	for err := range errorChan {
		if errors.Is(err, fs.ErrStreaming) {
			return []error{err}
		}
		retErrors = append(retErrors, err)
	}

	if len(retErrors) > 0 {
		return retErrors
	}

	site, err := b.Builder.Dispatch()
	if err != nil {
		return []error{err}
	}

	for _, plugin := range b.Plugins {
		if err := plugin.PreWrite(&site); err != nil {
			return []error{err}
		}
	}

	if err := b.Writer.Write(site); err != nil {
		return []error{err}
	}

	for _, plugin := range b.Plugins {
		if err := plugin.PostWrite(); err != nil {
			return []error{err}
		}
	}

	return nil
}

func (b *Build) processFile(contentDir, file string) error {
	src, err := ioutil.ReadFile(filepath.Join(contentDir, file))
	if err != nil {
		return err
	}

	page, err := b.Parser.ParsePage(src)
	if err != nil {
		return err
	}

	// For a file path like /blog/coffee/making-espresso.md, the resulting
	// page route will be /blog/coffee.
	page.Route = filepath.ToSlash(filepath.Dir(file))

	// For a file name like making-espresso.md, the resulting page ID will
	// be making-espresso.
	page.ID = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

	page.Href = filepath.Join(page.Route, page.ID)

	if err := b.setTypeForPage(&page); err != nil {
		return err
	}

	if err := b.Builder.RegisterPage(page); err != nil {
		return err
	}

	for _, plugin := range b.Plugins {
		if err := plugin.ProcessPage(&page); err != nil {
			return err
		}
	}

	return nil
}

// setTypeForPage sets the Type field of a page if a page type has been
// provided by the user. Returns an error if the provided page type has
// not been configured in the given types map.
func (b *Build) setTypeForPage(page *model.Page) error {
	providedType := page.ProvidedType()

	if providedType == "" {
		return nil
	}

	if _, exists := b.Types[providedType]; !exists {
		return fmt.Errorf("%s: type %s has not been declared in verless.yml", page.ID, providedType)
	}

	page.Type = b.Types[providedType]
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
