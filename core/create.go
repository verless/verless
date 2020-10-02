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

	if path != "." {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	} else {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			// RemoveAll removes nested directory in first iteration which causes
			// os.PathError saying "no such file or directory" for next recursion of
			// WalkFunc.
			if os.IsNotExist(err) {
				return nil
			}
			if path != "." {
				if info.IsDir() {
					// Remove nested non-empty directories as os.Remove() only removes
					// files and empty directories
					return os.RemoveAll(path)
				} else {
					return os.Remove(path)
				}
			}
			return nil
		})
		if err != nil {
			return errors.New("Cannot remove existing files from current directory.")
		}
	}

	dirs := []string{
		filepath.Join(path, ContentDir),
		filepath.Join(path, ThemesDir, DefaultTheme, TemplateDir),
		filepath.Join(path, ThemesDir, DefaultTheme, CSSDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	files := map[string][]byte{
		filepath.Join(path, "verless.yml"):                                     []byte(defaultConfig),
		filepath.Join(path, ".gitignore"):                                      []byte(defaultGitignore),
		filepath.Join(path, ThemesDir, DefaultTheme, TemplateDir, ListPageTpl): []byte(defaultTpl),
		filepath.Join(path, ThemesDir, DefaultTheme, TemplateDir, PageTpl):     {},
		filepath.Join(path, ThemesDir, DefaultTheme, CSSDir, "style.css"):      []byte(defaultCss),
	}

	return createFiles(files)
}

// CreateThemeOptions represents project path for creating new theme.
type CreateThemeOptions struct {
	Project string
}

// CreateTheme creates a new theme with the specified name inside the
// given path. Returns an error if it already exists, unless --overwrite
// has been used.
func CreateTheme(options CreateThemeOptions, name string) error {
	if _, err := os.Stat(options.Project); os.IsNotExist(err) {
		return ErrProjectNotExists
	}

	if _, err := os.Stat(filepath.Join(options.Project, ThemesDir, name)); !os.IsNotExist(err) {
		return ErrThemeExists
	}

	dirs := []string{
		filepath.Join(options.Project, ThemesDir, name, TemplateDir),
		filepath.Join(options.Project, ThemesDir, name, CSSDir),
		filepath.Join(options.Project, ThemesDir, name, JSDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	files := map[string][]byte{
		filepath.Join(options.Project, ThemesDir, name, TemplateDir, ListPageTpl): {},
		filepath.Join(options.Project, ThemesDir, name, TemplateDir, PageTpl):     {},
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
