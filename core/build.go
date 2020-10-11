package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/verless/verless/builder"
	"github.com/verless/verless/config"
	"github.com/verless/verless/fs"
	"github.com/verless/verless/model"
	"github.com/verless/verless/parser"
	"github.com/verless/verless/plugin"
	"github.com/verless/verless/theme"
	"github.com/verless/verless/writer"
)

const (
	// parallelism specifies the number of parallel workers.
	parallelism int = 4
)

var (
	// ErrCannotOverwrite states that verless isn't allowed to delete or
	// overwrite the output directory.
	ErrCannotOverwrite = errors.New(`cannot overwrite the output directory.
Consider using the --overwrite flag or enable build.overwrite in verless.yml`)

	// ErrMissingVersionKey states that the top-level `version` key is
	// empty or missing in verless.yml.
	ErrMissingVersionKey = errors.New("missing `version` key in verless.yml")
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
	Plugins []plugin.Plugin
	Types   map[string]*model.Type
	Options BuildOptions
}

// New initializes a new Build instance.
func NewBuild(targetFs afero.Fs, path string, options BuildOptions) (*Build, error) {
	cfg, err := config.FromFile(path, config.Filename)
	if err != nil {
		return nil, err
	}

	if cfg.Version == "" {
		return nil, ErrMissingVersionKey
	}

	outputDir := outputDir(path, &options)

	if !fs.IsSafeToRemove(targetFs, outputDir, options.Overwrite || cfg.Build.Overwrite) {
		return nil, ErrCannotOverwrite
	}

	writerCtx := writer.Context{
		Fs:                 targetFs,
		Path:               path,
		OutputDir:          outputDir,
		Theme:              cfg.Theme,
		RecompileTemplates: options.RecompileTemplates,
	}

	b := Build{
		Path:    path,
		Parser:  parser.NewMarkdown(),
		Builder: builder.New(&cfg),
		Writer:  writer.New(writerCtx),
		Types:   cfg.Types,
		Options: options,
	}

	plugins := plugin.LoadAll(&cfg, targetFs, outputDir)

	for _, key := range cfg.Plugins {
		if _, exists := plugins[key]; !exists {
			return nil, fmt.Errorf("plugin %s not found", key)
		}
		b.Plugins = append(b.Plugins, plugins[key]())
	}

	for _, beforeHook := range cfg.Build.Before {
		cmdParts := strings.Split(beforeHook, " ")
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		cmd.Dir = path
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}

	if err := theme.RunBeforeHooks(path, cfg.Theme); err != nil {
		return nil, err
	}

	return &b, nil
}

// Run executes the build using the provided build context.
//
// The current build implementation runs the following steps:
//	1. Read all files in the content directory and send them through a channel.
//	2. Spawn workers reading from that channel.
//	3. Process each received file:
//		3.1. Read the file as a []byte
//		3.2. Parse the []byte and convert it to a model.Page.
//		3.3. Register the page in the builder's site model.
//		3.4. Let each plugin process the page.
//	4. Get the site model from the builder and render it as a website.
//	5. Let each plugin finish its work, e.g. by writing a file.
func (b *Build) Run() error {
	var (
		files           = make(chan string)
		errorCh         = make(chan error)
		collectedErrors = make([]error, 0)
		contentDir      = filepath.Join(b.Path, config.ContentDir)
	)

	go func() {
		if err := fs.StreamFiles(contentDir, files, fs.MarkdownOnly, fs.NoUnderscores); err != nil {
			errorCh <- err
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(parallelism)

	for i := 0; i < parallelism; i++ {
		go func() {
			// Process the files received via the files channel.
			for file := range files {
				if err := b.processFile(contentDir, file); err != nil {
					errorCh <- err
				}
			}
			wg.Done()
		}()
	}

	// Close the error channel as soon as all workers have finished.
	go func() {
		for {
			wg.Wait()
			close(errorCh)
			break
		}
	}()

	for err := range errorCh {
		if errors.Is(err, fs.ErrStreaming) {
			return err
		}
		collectedErrors = append(collectedErrors, err)
	}

	if len(collectedErrors) > 0 {
		return fmt.Errorf("errors while processing files: %v", collectedErrors)
	}

	site, err := b.Builder.Dispatch()
	if err != nil {
		return err
	}

	for _, p := range b.Plugins {
		if err := p.PreWrite(&site); err != nil {
			return err
		}
	}

	if err := b.Writer.Write(site); err != nil {
		return err
	}

	for _, p := range b.Plugins {
		if err := p.PostWrite(); err != nil {
			return err
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

	// A page like /blog/coffee/making-espresso.md will have /blog/coffee as
	// route and making-espresso as ID.
	page.Route = filepath.ToSlash(filepath.Dir(file))
	page.ID = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	page.Href = filepath.ToSlash(filepath.Join(page.Route, page.ID))

	if err := b.setPageType(&page); err != nil {
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

// setPageType sets the Type field of a page if a page type has been
// provided by the user.
func (b *Build) setPageType(page *model.Page) error {
	providedType := page.ProvidedType()

	if providedType == "" {
		return nil
	}

	if _, exists := b.Types[providedType]; !exists {
		return fmt.
			Errorf("%s: type %s has not been declared", page.ID, providedType)
	}
	page.Type = b.Types[providedType]

	return nil
}

func outputDir(path string, options *BuildOptions) string {
	if options.OutputDir != "" {
		return options.OutputDir
	}

	return filepath.Join(path, config.OutputDir)
}
