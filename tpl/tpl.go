package tpl

import (
	"errors"
	"fmt"
	"text/template"
)

var (
	templates            map[string]*template.Template
	ErrAlreadyRegistered = errors.New("template has already been registered")
)

func Register(key string, path string, recompileTemplates bool) (*template.Template, error) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	if !recompileTemplates {
		if _, exists := templates[key]; exists {
			return nil, ErrAlreadyRegistered
		}
	}

	tpl, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}

	templates[key] = tpl

	return templates[key], nil
}

func Get(key string) (*template.Template, error) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	if _, exists := templates[key]; !exists {
		return nil, fmt.Errorf("template %s has not been registered", key)
	}

	return templates[key], nil
}

func IsRegistered(key string) bool {
	_, exists := templates[key]
	return exists
}
