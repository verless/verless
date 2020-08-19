package tpl

import (
	"github.com/verless/verless/config"
	"github.com/verless/verless/test"
	"path/filepath"
	"testing"
)

const (
	projectFolderPath = "../example"
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
	indexPageTplPath := filepath.Join(projectFolderPath, config.TemplateDir, config.IndexPageTpl)

	_, err := Register(config.IndexPageTpl, indexPageTplPath)
	test.Ok(t, err)

	_, err = Get(config.IndexPageTpl)
	test.Ok(t, err)

	_, err = Get(invalidKey)
	test.Assert(t, err != nil, "key is invalid")
}
