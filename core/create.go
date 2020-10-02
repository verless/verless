package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

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

	// ErrFileExist states that the specified file already exists.
	ErrFileExists = errors.New("file already exists.")

	// ErrNoSuchDirExists states that the specified directory doesn't exist
	ErrNoSuchDirExists = errors.New("no such directory exist, create it first")
)

// CreateProjectOptions represents options for creating a project.
type CreateProjectOptions struct {
	Overwrite bool
}

// CreateFileOptions represents project path for creating file.
type CreateFileOptions struct {
	Project string
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

// CreateTheme creates a new theme with the specified name inside the
// given path. Returns an error if it already exists, unless --overwrite
// has been used.
func CreateTheme(path, name string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrProjectNotExists
	}

	if _, err := os.Stat(filepath.Join(path, ThemesDir, name)); !os.IsNotExist(err) {
		return ErrThemeExists
	}

	dirs := []string{
		filepath.Join(path, ThemesDir, name, TemplateDir),
		filepath.Join(path, ThemesDir, name, CSSDir),
		filepath.Join(path, ThemesDir, name, JSDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	files := map[string][]byte{
		filepath.Join(path, ThemesDir, name, TemplateDir, ListPageTpl): {},
		filepath.Join(path, ThemesDir, name, TemplateDir, PageTpl):     {},
	}

	return createFiles(files)
}

// CreateFile creates a file with specified path under content directory.
func CreateFile(filePath string, options CreateFileOptions) error {

	if _, err := os.Stat(options.Project); os.IsNotExist(err) {
		return ErrProjectNotExists
	}

	contentPath := filepath.Join(options.Project, ContentDir, filePath)

	if _, err := os.Stat(path.Dir(contentPath)); os.IsNotExist(err) {
		return ErrNoSuchDirExists
	}

	if _, err := os.Stat(contentPath); !os.IsNotExist(err) {
		return ErrFileExists
	}

	defaultContentTemplate :=
		`---
Title:
Description:
Date: %s
---
`

	content := fmt.Sprintf(defaultContentTemplate, time.Now().Format("2006-01-02"))
	if err := ioutil.WriteFile(contentPath, []byte(content), 0755); err != nil {
		return err
	}

	return nil

}

func createFiles(files map[string][]byte) error {
	for path, content := range files {
		if err := ioutil.WriteFile(path, content, 0755); err != nil {
			return err
		}
	}
	return nil
}
