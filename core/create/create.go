package create

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/verless/verless/config"
	"github.com/verless/verless/fs"
)

// Project creates a new verless default project.
func Project(path string) error {
	if err := fs.MkdirAll(path, ContentDir,
		filepath.Join(ThemesDir, DefaultTheme, TemplateDir),
		filepath.Join(ThemesDir, DefaultTheme, "css"),
	); err != nil {
		return err
	}

	files := map[string][]byte{
		filepath.Join(path, "verless.yml"):                                     []byte(defaultConfig),
		filepath.Join(path, ThemesDir, DefaultTheme, TemplateDir, ListPageTpl): []byte(defaultTpl),
		filepath.Join(path, ThemesDir, DefaultTheme, TemplateDir, PageTpl):     {},
		filepath.Join(path, ThemesDir, DefaultTheme, CSSDir, "style.css"):      []byte(defaultCss),
	}

	return createFiles(files)
}

// Theme creates a new verless theme.
func Theme(path, name string) error {
	if err := fs.MkdirAll(path,
		filepath.Join(ThemesDir, name, TemplateDir),
		filepath.Join(ThemesDir, name, "css"),
		filepath.Join(ThemesDir, name, "js"),
	); err != nil {
		return err
	}

	files := map[string][]byte{
		filepath.Join(path, ThemesDir, name, TemplateDir, ListPageTpl): {},
		filepath.Join(path, ThemesDir, name, TemplateDir, PageTpl):     {},
	}

	return createFiles(files)
}

// createFiles takes a map of file paths mapped against the
// file contents, creates those file paths and writes the contents
// into them.
//
// The keys have to be file paths like `my-blog/verless.yml` and
// all directories already have to exist.
func createFiles(files map[string][]byte) error {
	for path, content := range files {
		if err := ioutil.WriteFile(path, content, 0755); err != nil {
			return err
		}
	}

	return nil
}
