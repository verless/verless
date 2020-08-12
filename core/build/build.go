// Package build provides verless' core build functionality and
// an entry point for it. All public functions may only be called
// by the outer core package functions.
package build

import (
	"io/ioutil"
	"path/filepath"

	"github.com/verless/verless/config"
	"github.com/verless/verless/fs"
	"github.com/verless/verless/model"
)

const (
	// parallelism specifies the number of parallel workers.
	parallelism = 4
)

type (
	// Parser represents a parser that processes Markdown files
	// and converts them into a model instance.
	Parser interface {
		// ParsePage must be safe for concurrent usage.
		ParsePage(src []byte) (model.Page, error)
	}

	// Builder represents a model builder that maintains a Site
	// instance and registers all parsed pages in that instance.
	Builder interface {
		// RegisterPage must be safe for concurrent usage.
		RegisterPage(route string, page model.Page) error
		Dispatch() (model.Site, error)
	}

	// Writer represents a model writer that renders the site
	// model as HTML using the corresponding templates.
	Writer interface {
		Write(site model.Site) error
	}

	// Plugin represents a built-in verless plugin.
	Plugin interface {
		// ProcessPage will be invoked after parsing the page.
		// Must be safe for concurrent usage.
		ProcessPage(route string, page *model.Page) error
		// Finalize will be invoked after rendering the site.
		Finalize() error
	}
)

// Context provides all components required for running a build.
type Context struct {
	Path    string
	Parser  Parser
	Builder Builder
	Writer  Writer
	Plugins []Plugin
}

// Run executes the build using the provided build context.
//
// The current build implementation runs the following steps to
// build the static site:
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
// For further info on one of these steps, see its implementation.
func Run(ctx Context) error {
	var (
		files      = make(chan string)
		errors     = make(chan error)
		contentDir = filepath.Join(ctx.Path, config.ContentDir)
	)

	go func() {
		errors <- fs.StreamFiles(contentDir, files, fs.MarkdownOnly, fs.NoUnderscores)
	}()

	err := processFiles(func(file string) error {
		src, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		page, err := ctx.Parser.ParsePage(src)
		if err != nil {
			return err
		}

		// Given a file path like `/blog/coffee/making-espresso.md`,
		// the resulting path for the page will be `/blog/coffee`.
		path := filepath.Dir(file)

		if err := ctx.Builder.RegisterPage(path, page); err != nil {
			return err
		}

		for _, plugin := range ctx.Plugins {
			if err := plugin.ProcessPage(path, &page); err != nil {
				return err
			}
		}

		return nil
	}, files, parallelism)

	if err != nil {
		return err
	}

	site, err := ctx.Builder.Dispatch()
	if err != nil {
		return err
	}

	if err := ctx.Writer.Write(site); err != nil {
		return err
	}

	for _, plugin := range ctx.Plugins {
		// Finalize has to be called after writing the page to make
		// sure that no files will be overwritten by the Writer.
		if err := plugin.Finalize(); err != nil {
			return err
		}
	}

	return nil
}
