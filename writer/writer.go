package writer

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/otiai10/copy"
	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
)

func New(path, outputDir string) (*writer, error) {
	w := writer{
		path:       path,
		outputPath: outputDir,
	}

	if err := w.initTemplates(); err != nil {
		return nil, err
	}

	return &w, nil
}

type writer struct {
	path         string
	outputPath   string
	site         model.Site
	pageTpl      *template.Template
	indexPageTpl *template.Template
}

func (w *writer) Write(site model.Site) error {
	w.site = site

	err := w.site.WalkRoutes(func(path string, route *model.Route) error {
		for _, p := range route.Pages {
			if err := w.writePage(path, page{
				Meta:   &w.site.Meta,
				Nav:    &w.site.Nav,
				Page:   &p,
				Footer: &w.site.Footer,
			}); err != nil {
				return err
			}
		}
		return w.writeIndexPage(path, indexPage{
			Meta:   &w.site.Meta,
			Nav:    &w.site.Nav,
			Page:   &route.IndexPage,
			Footer: &w.site.Footer,
		})
	}, -1)

	if err != nil {
		return err
	}

	if err := w.copyAssetDir(); err != nil {
		return err
	}

	return nil
}

func (w *writer) initTemplates() error {
	var (
		err          error
		pageTpl      = filepath.Join(w.path, config.TemplateDir, config.PageTpl)
		indexPageTpl = filepath.Join(w.path, config.TemplateDir, config.IndexPageTpl)
	)

	if w.pageTpl, err = template.ParseFiles(pageTpl); err != nil {
		return err
	}

	if w.indexPageTpl, err = template.ParseFiles(indexPageTpl); err != nil {
		return err
	}

	return nil
}

func (w *writer) writePage(route string, page page) error {
	path := filepath.Join(w.outputPath, route, page.Page.ID)

	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, config.IndexFile))
	if err != nil {
		return err
	}

	return w.pageTpl.Execute(file, &page)
}

func (w *writer) writeIndexPage(route string, indexPage indexPage) error {
	path := filepath.Join(w.outputPath, route)

	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, config.IndexFile))
	if err != nil {
		return err
	}

	return w.indexPageTpl.Execute(file, &indexPage)
}

func (w *writer) copyAssetDir() error {
	var (
		src  = filepath.Join(w.path, config.AssetDir)
		dest = filepath.Join(w.outputPath, config.AssetDir)
	)

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}

	return copy.Copy(src, dest)
}
