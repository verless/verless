package tags

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
)

const (
	Key     string = "tags"
	tagsDir string = "tags"
)

func New(path, outputDir string) *tags {
	t := tags{
		path:      path,
		outputDir: outputDir,
		m:         make(map[string]*model.IndexPage),
	}

	return &t
}

type tags struct {
	path      string
	outputDir string
	m         map[string]*model.IndexPage
}

func (t *tags) ProcessPage(_ string, page *model.Page) error {
	for _, tag := range page.Tags {
		if _, exists := t.m[tag]; !exists {
			t.createIndexPage(tag)
		}
		t.m[tag].Pages = append(t.m[tag].Pages, page)
	}

	return nil
}

func (t *tags) Finalize(site *model.Site) error {
	tpl, err := template.ParseFiles(filepath.Join(t.path, config.TemplateDir, config.IndexPageTpl))
	if err != nil {
		return err
	}

	for tag, ip := range t.m {
		if err := t.writeIndexPage(tag, ip, tpl, site); err != nil {
			return err
		}
	}

	return nil
}

func (t *tags) createIndexPage(key string) {
	t.m[key] = &model.IndexPage{
		Pages: make([]*model.Page, 0),
	}
}

func (t *tags) writeIndexPage(tag string, ip *model.IndexPage, tpl *template.Template, site *model.Site) error {

	path := filepath.Join(t.outputDir, tagsDir, strings.ToLower(tag))
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, config.IndexFile))
	if err != nil {
		return err
	}

	indexPage := indexPage{
		Meta:      &site.Meta,
		Nav:       &site.Nav,
		IndexPage: ip,
		Footer:    &site.Footer,
	}

	return tpl.Execute(file, indexPage)
}
