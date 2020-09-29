package core

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	. "github.com/verless/verless/config"
	"github.com/verless/verless/fs"
)

var (
	// ErrProjectExists states that the specified project already exists.
	ErrProjectExists = errors.New("project already exists, use --overwrite to remove it")

	// ErrProjectNotExists states that the specified project doesn't exist.
	ErrProjectNotExists = errors.New("project doesn't exist yet, create it first")

	// ErrThemeExists states that the specified theme already exists.
	ErrThemeExists = errors.New("theme already exists, remove it first")
)

// CreateProjectOptions represents options for creating a project.
type CreateProjectOptions struct {
	Overwrite bool
}

// CreateProject creates a new verless project. If the specified project
// path already exists, CreateProject returns an error unless --overwrite
// has been used.
func CreateProject(path string, options CreateProjectOptions) error {
	if !fs.IsSafeToRemove(afero.NewOsFs(), path, options.Overwrite) {
		return ErrProjectExists
	}

	if err := os.RemoveAll(path); err != nil {
		return err
	}

	dirs := []string{
		filepath.Join(path, ContentDir),
		theme.TemplateDir(path, DefaultTheme),
		theme.CssDir(path, DefaultTheme),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	files := map[string][]byte{
		filepath.Join(path, "verless.yml"):                                defaultConfig,
		filepath.Join(theme.TemplateDir(path, DefaultTheme), ListPageTpl): defaultTpl,
		filepath.Join(theme.TemplateDir(path, DefaultTheme), PageTpl):     {},
		filepath.Join(theme.CssDir(path, DefaultTheme), "style.css"):      defaultCss,
	}

	return createFiles(files)
}

// CreateTheme creates a new theme with the specified name inside the
// given path. Returns an error if it already exists, unless --overwrite
// has been used.
func CreateTheme(path, name string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrProjectNotExists
	}

	if theme.Exists(path, name) {
		return ErrThemeExists
	}

	dirs := []string{
		theme.TemplateDir(path, name),
		theme.CssDir(path, name),
		theme.JsDir(path, name),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	files := map[string][]byte{
		filepath.Join(theme.TemplateDir(path, name), ListPageTpl): {},
		filepath.Join(theme.TemplateDir(path, name), PageTpl):     {},
		filepath.Join(theme.Dir(path, name), "theme.yml"):         defaultThemeConfig,
	}

	return createFiles(files)
}

func createFiles(files map[string][]byte) error {
	for path, content := range files {
		if err := ioutil.WriteFile(path, content, 0755); err != nil {
			return err
		}
	}
	return nil
}
