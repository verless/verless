// Package tpl provides a simple template registry and corresponding
// function for registering and loading templates from that registry.
package tpl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

var (
	// templates stores all successfully rendered template instances.
	templates map[string]*template.Template

	// ErrAlreadyRegistered is returned when a template with a given
	// key has already been registered.
	ErrAlreadyRegistered = errors.New("template has already been registered")
)

// Register parses a template file and registers the instance under
// the given key. If a template with the key has already registered,
// Register will return an error unless the registration is forced.
func Register(key string, path string, basePath string, force bool) (*template.Template, error) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	if !force {
		if _, exists := templates[key]; exists {
			return nil, ErrAlreadyRegistered
		}
	}

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		tpl, err := template.ParseFiles(path)
		if err != nil {
			return nil, err
		}

		templates[key] = tpl
		return templates[key], nil
	}

	return registerWithBaseTemplates(key, path, basePath)
}

func registerWithBaseTemplates(key string, path string, basePath string) (*template.Template, error) {
	baseFiles, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	fileNames := make([]string, len(baseFiles)+1)
	fileNames[0] = path
	for i, f := range baseFiles {
		fileNames[i+1] = filepath.Join(basePath, f.Name())
	}

	tpl, err := template.ParseFiles(fileNames...)
	templates[key] = tpl

	return templates[key], nil
}

// Get returns the template registered under the given key.
//
// The template must be registered using Register first. If it hasn't
// been registered successfully before, Get will return an error. Use
// IsRegistered if you're unsure if the template exists.
func Get(key string) (*template.Template, error) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	if _, exists := templates[key]; !exists {
		return nil, fmt.Errorf("template %s has not been registered", key)
	}

	return templates[key], nil
}

// IsRegistered indicates whether a template with the given key has
// been registered successfully using Register.
func IsRegistered(key string) bool {
	_, exists := templates[key]
	return exists
}
