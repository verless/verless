package tpl

import (
	"path/filepath"
	"testing"
	"text/template"

	"github.com/verless/verless/config"
	"github.com/verless/verless/test"
)

const (
	projectFolderPath = "../example"
	testKey           = "test key"
	invalidKey        = "invalid key"
)

func TestRegister(t *testing.T) {
	pageTplPath := filepath.Join(projectFolderPath, config.TemplateDir, config.PageTpl)

	_, err := Register(config.PageTpl, pageTplPath)
	test.Ok(t, err)

	_, err = Register(config.PageTpl, pageTplPath)
	test.Assert(t, err != nil, "template has already been registered")
}

func TestGet(t *testing.T) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates[testKey] = &template.Template{}

	tpl, err := Get(testKey)
	test.Ok(t, err)

	test.Assert(t, tpl == templates[testKey], "template has to be in map")

	_, err = Get(invalidKey)
	test.Assert(t, err != nil, "template key is invalid")
}
