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
	"github.com/verless/verless/theme"
)

var (
	// ErrProjectExists states that the specified project already exists.
	ErrProjectExists = errors.New("project already exists, use --overwrite to remove it")

	// ErrProjectNotExists states that the specified project doesn't exist.
	ErrProjectNotExists = errors.New("project doesn't exist yet, create it first")

	// ErrThemeExists states that the specified theme already exists.
	ErrThemeExists = errors.New("theme already exists, remove it first")

	// ErrFileExist states that the specified file already exists.
	ErrFileExists = errors.New("file already exists")
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
			return errors.New("Cannot remove existing files from current directory")
		}
	}

	dirs := []string{
		filepath.Join(path, ContentDir),
		theme.TemplatePath(path, theme.Default),
		theme.CssPath(path, theme.Default),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	files := map[string][]byte{
		filepath.Join(path, "verless.yml"):                                             defaultConfig,
		filepath.Join(path, ".gitignore"):                                              defaultGitignore,
		filepath.Join(theme.TemplatePath(path, theme.Default), theme.ListPageTemplate): defaultTpl,
		filepath.Join(theme.TemplatePath(path, theme.Default), theme.PageTemplate):     {},
		filepath.Join(theme.CssPath(path, theme.Default), "style.css"):                 defaultCss,
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

	if theme.Exists(options.Project, name) {
		return ErrThemeExists
	}

	dirs := []string{
		theme.TemplatePath(options.Project, name),
		theme.CssPath(options.Project, name),
		theme.JsPath(options.Project, name),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	files := map[string][]byte{
		filepath.Join(theme.TemplatePath(options.Project, name), theme.ListPageTemplate): {},
		filepath.Join(theme.TemplatePath(options.Project, name), theme.PageTemplate):     {},
		filepath.Join(theme.Path(options.Project, name), "theme.yml"):                    defaultThemeConfig,
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
		return fmt.Errorf("no such dir %s exist, create it first", path.Dir(contentPath))
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
