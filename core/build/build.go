// Package build provides verless' core build functionality.
package build

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

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
		RegisterPage(page model.Page) error
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
		ProcessPage(page *model.Page) error
		// PreWrite will be invoked before writing the site.
		PreWrite(site *model.Site) error
		// PostWrite will be invoked after writing the site.
		PostWrite() error
	}
)

// Context provides all components required for running a build.
type Context struct {
	Path               string
	Parser             Parser
	Builder            Builder
	Writer             Writer
	Plugins            []Plugin
	Types              map[string]*model.Type
	RecompileTemplates bool
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
//
// If any error occurs it is returned as a slice of errors as there can be
// several errors from the concurrent goroutines at the same time.
// If any error occurs all goroutines are stopped as soon as possible.
func Run(ctx Context) []error {
	var (
		files      = make(chan string)
		errorChan  = make(chan error)
		retErrors  = make([]error, 0)
		contentDir = filepath.Join(ctx.Path, config.ContentDir)
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
				if err := processFile(&ctx, contentDir, file); err != nil {
					errorChan <- err
				}
			}
			wg.Done()
		}()
	}

	// Observe the WaitGroup and close the error channel when
	// all workers have finished.
	go func() {
		for {
			wg.Wait()
			close(errorChan)
			break
		}
	}()

	// Collect all errors. This automatically stops when the
	// workers have finished, as the channel gets closed.
	for err := range errorChan {
		if errors.Is(err, fs.ErrStreaming) {
			return []error{err}
		}
		retErrors = append(retErrors, err)
	}

	if len(retErrors) > 0 {
		return retErrors
	}

	site, err := ctx.Builder.Dispatch()
	if err != nil {
		return []error{err}
	}

	for _, plugin := range ctx.Plugins {
		if err := plugin.PreWrite(&site); err != nil {
			return []error{err}
		}
	}

	if err := ctx.Writer.Write(site); err != nil {
		return []error{err}
	}

	for _, plugin := range ctx.Plugins {
		if err := plugin.PostWrite(); err != nil {
			return []error{err}
		}
	}

	return nil
}

func processFile(ctx *Context, contentDir, file string) error {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	page, err := ctx.Parser.ParsePage(src)
	if err != nil {
		return err
	}

	// For a file path like example/content/blog/coffee/making-espresso.md,
	// the resulting path will be /blog/coffee.
	path := filepath.Dir(file)[len(contentDir):]
	if path == "" {
		path = "/"
	}
	page.Route = filepath.ToSlash(path)

	// For a file name like making-espresso.md, the resulting page
	// ID will be making-espresso.
	page.ID = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

	if err := setType(&page, ctx.Types); err != nil {
		return err
	}

	if err := ctx.Builder.RegisterPage(page); err != nil {
		return err
	}

	for _, plugin := range ctx.Plugins {
		if err := plugin.ProcessPage(&page); err != nil {
			return err
		}
	}

	return nil
}

// setType sets the Type field of a page if a page type has been
// provided by the user. Returns an error if the provided page type
// has not been configured in the given types map.
func setType(page *model.Page, types map[string]*model.Type) error {
	providedType := page.ProvidedType()

	if providedType == "" {
		return nil
	}

	if _, exists := types[providedType]; !exists {
		return fmt.Errorf("%s: type %s has not been declared in verless.yml", page.ID, providedType)
	}

	page.Type = types[providedType]
	return nil
}
