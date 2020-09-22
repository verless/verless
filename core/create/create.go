package create

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/verless/verless/config"
	"github.com/verless/verless/fs"
)

// Project creates a new verless default project.
func Project(path string) error {
	err := fs.MkdirAll(path, ContentDir, filepath.Join(ThemesDir, DefaultTheme, TemplateDir),
		filepath.Join(ThemesDir, DefaultTheme, "css"))
	if err != nil {
		return err
	}

	files := map[string][]byte{
		filepath.Join(path, "verless.yml"):                                     []byte(defaultConfig),
		filepath.Join(path, ThemesDir, DefaultTheme, TemplateDir, ListPageTpl): []byte(defaultTpl),
		filepath.Join(path, ThemesDir, DefaultTheme, TemplateDir, PageTpl):     {},
		filepath.Join(path, ThemesDir, DefaultTheme, CSSDir, CSSFile):          []byte(defaultCss),
	}

	if err := createFiles(files); err != nil {
		return err
	}

	return nil
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
