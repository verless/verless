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
	"github.com/verless/verless/tpl"
	"github.com/verless/verless/tree"
)

type Context struct {
	Fs                 afero.Fs
	Path               string
	OutputDir          string
	Theme              string
	OverwriteTemplates bool
}

// New creates a new writer that renders the site model in the given
// filesystem instance to outputDir.
func New(ctx Context) *writer {
	if ctx.Theme == "" {
		ctx.Theme = config.DefaultTheme
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

	file, err := w.ctx.Fs.Create(filepath.Join(path, config.IndexFile))
	if err != nil {
		return err
	}

	pageTpl, err := w.loadTemplate(page.Page.Type, config.PageTpl)
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

	file, err := w.ctx.Fs.Create(filepath.Join(path, config.IndexFile))
	if err != nil {
		return err
	}

	listPageTpl, err := w.loadTemplate(listPage.Type, config.ListPageTpl)
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

	if !w.ctx.OverwriteTemplates && tpl.IsRegistered(pageTpl) {
		return tpl.Get(pageTpl)
	}

	tplPath := filepath.Join(w.ctx.Path, config.ThemesDir, w.ctx.Theme, config.TemplateDir, pageTpl)

	return tpl.Register(pageTpl, tplPath, w.ctx.OverwriteTemplates)
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
			src:      filepath.Join(w.ctx.Path, config.ThemesDir, w.ctx.Theme, config.CSSDir),
			dest:     filepath.Join(w.ctx.OutputDir, config.CSSDir),
			fileOnly: true,
		},
		{
			src:      filepath.Join(w.ctx.Path, config.ThemesDir, w.ctx.Theme, config.JSDir),
			dest:     filepath.Join(w.ctx.OutputDir, config.JSDir),
			fileOnly: true,
		},
	}

	for _, dir := range dirs {
		if err := fs.CopyFromOS(w.ctx.Fs, dir.src, dir.dest, dir.fileOnly); err != nil {
			return err
		}
	}

	return nil
}
