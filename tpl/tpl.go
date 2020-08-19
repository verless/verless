package tpl

import (
	"fmt"
	"text/template"
)

var (
	templates map[string]*template.Template
)

func Register(key string, path string) (*template.Template, error) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	if _, exists := templates[key]; !exists {
		tpl, err := template.ParseFiles(path)
		if err != nil {
			return nil, err
		}

		templates[key] = tpl
	}

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
