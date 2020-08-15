package writer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/otiai10/copy"
	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
)

func New(path, outputDir string) (*writer, error) {
	w := writer{
		path:      path,
		outputDir: outputDir,
	}

	if err := w.initTemplates(); err != nil {
		return nil, err
	}

	return &w, nil
}

type writer struct {
	path         string
	outputDir    string
	site         model.Site
	pageTpl      *template.Template
	indexPageTpl *template.Template
}

func (w *writer) Write(site model.Site) error {
	if err := w.removeOutDirIfExists(); err != nil {
		return err
	}

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
			Meta:      &w.site.Meta,
			Nav:       &w.site.Nav,
			IndexPage: &route.IndexPage,
			Footer:    &w.site.Footer,
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
	path := filepath.Join(w.outputDir, route, page.Page.ID)

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
	path := filepath.Join(w.outputDir, route)

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
		dest = filepath.Join(w.outputDir, config.AssetDir)
	)

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}

	return copy.Copy(src, dest)
}

// removeOutDirIfExists checks if the output folder exists and tries
// to remove it if it does.
func (w *writer) removeOutDirIfExists() error {
	info, err := ioutil.ReadDir(w.outputDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(info) > 0 {
		return os.RemoveAll(w.outputDir)
	}

	return nil
}
