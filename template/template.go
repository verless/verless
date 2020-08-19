package template

import "text/template"

var (
	templates map[string]*template.Template
)

func Load(path string) (*template.Template, error) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	if _, exists := templates[path]; !exists {
		tpl, err := template.ParseFiles(path)
		if err != nil {
			return nil, err
		}

		templates[path] = tpl
	}

	return templates[path], nil
}
