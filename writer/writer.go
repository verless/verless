// Package writer provides the default build.Writer implementation
// capable of writing the site model to a filesystem.
package writer

import (
	"path/filepath"
	"text/template"

	"github.com/spf13/afero"
	"github.com/verless/verless/config"
	"github.com/verless/verless/fs"
	"github.com/verless/verless/model"
	"github.com/verless/verless/theme"
	"github.com/verless/verless/tpl"
	"github.com/verless/verless/tree"
)

const (
	indexFile string = "index.html"
)

type Context struct {
	Fs                 afero.Fs
	Path               string
	OutputDir          string
	Theme              string
	RecompileTemplates bool
}

// New creates a new writer that renders the site model in the given
// filesystem instance to outputDir.
func New(ctx Context) *writer {
	if ctx.Theme == "" {
		ctx.Theme = theme.Default
	}

	w := writer{ctx: ctx}

	return &w
}

type writer struct {
	site model.Site
	ctx  Context
}

// Write renders the entire site model to the writer's filesystem.
//
// Basically, it creates a directory for each page and renders the
// page using its respective template. It also copies all assets.
func (w *writer) Write(site model.Site) error {
	if err := fs.Rmdir(w.ctx.Fs, w.ctx.OutputDir); err != nil {
		return err
	}

	w.site = site

	err := tree.Walk(w.site.Root, func(_ string, node tree.Node) error {
		for _, p := range node.(*model.Node).Pages {
			if err := w.writePage(p.Route, page{
				Meta:   &w.site.Meta,
				Nav:    &w.site.Nav,
				Page:   &p,
				Footer: &w.site.Footer,
			}); err != nil {
				return err
			}
		}

		lp := node.(*model.Node).ListPage

		if lp.Route == "" {
			panic("route must not be empty")
		}

		return w.writeListPage(lp.Route, listPage{
			Meta:     &w.site.Meta,
			Nav:      &w.site.Nav,
			ListPage: &lp,
			Footer:   &w.site.Footer,
		})
	}, -1)

	if err != nil {
		return err
	}

	if err := w.copyDirs(); err != nil {
		return err
	}

	return nil
}

// writePage renders a single page by applying the associated template
// and writing the file inside the output directory.
func (w *writer) writePage(route string, page page) error {
	path := filepath.Join(w.ctx.OutputDir, route, page.Page.ID)

	if err := w.ctx.Fs.MkdirAll(path, 0700); err != nil {
		return err
	}

	file, err := w.ctx.Fs.Create(filepath.Join(path, indexFile))
	if err != nil {
		return err
	}

	pageTpl, err := w.loadTemplate(page.Page.Type, theme.PageTemplate)
	if err != nil {
		return err
	}

	return pageTpl.Execute(file, &page)
}

// writeListPage does the same thing as writePage but for list pages.
func (w *writer) writeListPage(route string, listPage listPage) error {
	path := filepath.Join(w.ctx.OutputDir, route)

	if err := w.ctx.Fs.MkdirAll(path, 0700); err != nil {
		return err
	}

	file, err := w.ctx.Fs.Create(filepath.Join(path, indexFile))
	if err != nil {
		return err
	}

	listPageTpl, err := w.loadTemplate(listPage.Type, theme.ListPageTemplate)
	if err != nil {
		return err
	}

	return listPageTpl.Execute(file, &listPage)
}

// loadTemplate considers a page type and a default template, decides
// which template to use and loads that template from the registry.
func (w *writer) loadTemplate(t *model.Type, defaultTpl string) (*template.Template, error) {
	var pageTpl string

	switch {
	case t != nil && t.Template != "":
		pageTpl = t.Template
	default:
		pageTpl = defaultTpl
	}

	if !w.ctx.RecompileTemplates && tpl.IsRegistered(pageTpl) {
		return tpl.Get(pageTpl)
	}

	templatesRoot := filepath.Join(theme.TemplatePath(w.ctx.Path, w.ctx.Theme))
	tplPath := filepath.Join(templatesRoot, pageTpl)
	basePath := filepath.Join(templatesRoot, theme.TemplateBase)

	return tpl.Register(pageTpl, tplPath, basePath, w.ctx.RecompileTemplates)
}

func (w *writer) copyDirs() error {
	dirs := []struct {
		src      string
		dest     string
		fileOnly bool
	}{
		{
			src:      filepath.Join(w.ctx.Path, config.StaticDir),
			dest:     filepath.Join(w.ctx.OutputDir, config.StaticDir),
			fileOnly: false,
		},
		{
			src:      theme.CssPath(w.ctx.Path, w.ctx.Theme),
			dest:     filepath.Join(w.ctx.OutputDir, theme.CssDir),
			fileOnly: true,
		},
		{
			src:      theme.JsPath(w.ctx.Path, w.ctx.Theme),
			dest:     filepath.Join(w.ctx.OutputDir, theme.JsDir),
			fileOnly: true,
		},
		{
			src:      theme.AssetsPath(w.ctx.Path, w.ctx.Theme),
			dest:     filepath.Join(w.ctx.OutputDir, theme.AssetsDir),
			fileOnly: false,
		},
		{
			src:      theme.GeneratedPath(w.ctx.Path, w.ctx.Theme),
			dest:     filepath.Join(w.ctx.OutputDir, theme.GeneratedDir),
			fileOnly: false,
		},
	}

	for _, dir := range dirs {
		if err := fs.CopyFromOS(w.ctx.Fs, dir.src, dir.dest, dir.fileOnly); err != nil {
			return err
		}
	}

	return nil
}
