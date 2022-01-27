package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
	. "github.com/verless/verless/config"
	"github.com/verless/verless/fs"
	"github.com/verless/verless/theme"
)

const gitDirectory string = ".git"

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

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range dir {
		if strings.Contains(file.Name(), gitDirectory) {
			continue
		}
		if err := os.RemoveAll(filepath.Join(path, file.Name())); err != nil {
			return err
		}
	}

	dirs := []string{
		filepath.Join(path, ContentDir),
		filepath.Join(path, StaticDir),
		theme.TemplatePath(path, theme.Default),
		theme.AssetsPath(path, theme.Default),
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
		filepath.Join(theme.AssetsPath(path, theme.Default), "style.css"):              defaultCss,
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
		theme.AssetsPath(options.Project, name),
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
		createDirIfNotExist(path.Dir(contentPath))
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
	if err := ioutil.WriteFile(contentPath, []byte(content), 0644); err != nil {
		return err
	}

	return nil

}

func createFiles(files map[string][]byte) error {
	for path, content := range files {
		if err := ioutil.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}
	return nil
}

func createDirIfNotExist(filePath string) error {
	if err := os.Mkdir(filePath, 0755); err != nil {
		return err
	}
	return nil
}
