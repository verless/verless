// Package writer provides the default build.Writer implementation
// capable of writing the site model to a filesystem.
package writer

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/otiai10/copy"
	"github.com/spf13/afero"
	"github.com/verless/verless/config"
	"github.com/verless/verless/fs"
	"github.com/verless/verless/model"
	"github.com/verless/verless/tpl"
	"github.com/verless/verless/tree"
)

// New creates a new writer that renders the site model in the given
// filesystem instance to outputDir.
func New(fs afero.Fs, path, theme, outputDir string, recompileTemplates bool) *writer {
	if theme == "" {
		theme = config.DefaultTheme
	}

	w := writer{
		fs:                 fs,
		path:               path,
		theme:              theme,
		outputDir:          outputDir,
		recompileTemplates: recompileTemplates,
	}

	return &w
}

type writer struct {
	fs                 afero.Fs
	path               string
	theme              string
	outputDir          string
	site               model.Site
	recompileTemplates bool
}

// Write renders the entire site model to the writer's filesystem.
//
// Basically, it creates a directory for each page and renders the
// page using its respective template. It also copies all assets.
func (w *writer) Write(site model.Site) error {
	if err := fs.Rmdir(w.fs, w.outputDir); err != nil {
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

	if err := w.copyStaticDir(); err != nil {
		return err
	}

	if err := w.copyThemeDirs(); err != nil {
		return err
	}

	return nil
}

// writePage renders a single page by applying the associated template
// and writing the file inside the output directory.
func (w *writer) writePage(route string, page page) error {
	path := filepath.Join(w.outputDir, route, page.Page.ID)

	if err := w.fs.MkdirAll(path, 0700); err != nil {
		return err
	}

	file, err := w.fs.Create(filepath.Join(path, config.IndexFile))
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
	path := filepath.Join(w.outputDir, route)

	if err := w.fs.MkdirAll(path, 0700); err != nil {
		return err
	}

	file, err := w.fs.Create(filepath.Join(path, config.IndexFile))
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

	if !w.recompileTemplates && tpl.IsRegistered(pageTpl) {
		return tpl.Get(pageTpl)
	}

	tplPath := filepath.Join(w.path, config.ThemesDir, w.theme, config.TemplateDir, pageTpl)

	return tpl.Register(pageTpl, tplPath, w.recompileTemplates)
}

func (w *writer) copyStaticDir() error {
	// If the writer's target filesystem is the OS filesystem, directly
	// copy the asset directory using the copy package.
	if _, ok := w.fs.(*afero.OsFs); ok {
		var (
			src  = filepath.Join(w.path, config.StaticDir)
			dest = filepath.Join(w.outputDir, config.StaticDir)
		)
		if _, err := w.fs.Stat(src); os.IsNotExist(err) {
			return nil
		}
		return copy.Copy(src, dest)
	} else {
		// Otherwise, copy the assets directory from the physical filesystem
		// into the memory filesystem.
		return fs.CopyFromOS(w.fs, w.path, w.outputDir, false)
	}
}

func (w *writer) copyThemeDirs() error {
	dirs := []struct {
		src  string
		dest string
	}{
		{
			src:  filepath.Join(w.path, config.ThemesDir, w.theme, config.CSSDir),
			dest: filepath.Join(w.outputDir, config.CSSDir),
		},
		{
			src:  filepath.Join(w.path, config.ThemesDir, w.theme, config.JSDir),
			dest: filepath.Join(w.outputDir, config.JSDir),
		},
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir.src); os.IsNotExist(err) {
			continue
		}
		if err := fs.CopyFromOS(w.fs, dir.src, dir.dest, true); err != nil {
			return err
		}
	}

	return nil
}
