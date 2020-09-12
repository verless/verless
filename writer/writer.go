package writer

import (
	"io/ioutil"
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

func New(fs afero.Fs, path, outputDir string, recompileTemplates bool) *writer {
	w := writer{
		fs:                 fs,
		path:               path,
		outputDir:          outputDir,
		recompileTemplates: recompileTemplates,
	}

	return &w
}

type writer struct {
	fs                 afero.Fs
	path               string
	outputDir          string
	site               model.Site
	recompileTemplates bool
}

func (w *writer) Write(site model.Site) error {
	if err := fs.Rmdir(w.fs, w.outputDir); err != nil {
		return err
	}

	w.site = site

	err := tree.Walk(w.site.Root, func(node tree.Node) error {
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

		ip := node.(*model.Node).IndexPage

		return w.writeIndexPage(ip.Route, indexPage{
			Meta:      &w.site.Meta,
			Nav:       &w.site.Nav,
			IndexPage: &ip,
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

func (w *writer) writeIndexPage(route string, indexPage indexPage) error {
	path := filepath.Join(w.outputDir, route)

	if err := w.fs.MkdirAll(path, 0700); err != nil {
		return err
	}

	file, err := w.fs.Create(filepath.Join(path, config.IndexFile))
	if err != nil {
		return err
	}

	indexPageTpl, err := w.loadTemplate(indexPage.Type, config.IndexPageTpl)
	if err != nil {
		return err
	}

	return indexPageTpl.Execute(file, &indexPage)
}

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

	tplPath := filepath.Join(w.path, config.TemplateDir, pageTpl)

	return tpl.Register(pageTpl, tplPath, w.recompileTemplates)
}

func (w *writer) copyAssetDir() error {
	// If the writer's target filesystem is the OS filesystem, directly
	// copy the asset directory using the copy package.
	if _, ok := w.fs.(*afero.OsFs); ok {
		var (
			src  = filepath.Join(w.path, config.AssetDir)
			dest = filepath.Join(w.outputDir, config.AssetDir)
		)
		if _, err := w.fs.Stat(src); os.IsNotExist(err) {
			return nil
		}
		return copy.Copy(src, dest)
	} else {
		return copyFromOsFs(w.fs, w.path, w.outputDir)
	}
}

func copyFromOsFs(targetFs afero.Fs, src, dest string) error {
	var (
		files = make(chan string)
		err   error
	)

	go func() {
		err = fs.StreamFiles(src, files)
	}()

	for file := range files {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		path := file[len(src):]
		path = filepath.Join(dest, path)

		memFile, err := targetFs.Create(path)
		if err != nil {
			return err
		}
		if _, err := memFile.Write(bytes); err != nil {
			return err
		}
	}
	return err
}
