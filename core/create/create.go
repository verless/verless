package create

import (
	"io/ioutil"
	"path/filepath"

	"github.com/verless/verless/config"
	"github.com/verless/verless/fs"
)

// Project creates a new verless project which is an exact copy
// of the project at the given path.
func Project(path string) error {
	if err := fs.MkdirAll(path, config.AssetDir, config.ContentDir, config.TemplateDir); err != nil {
		return err
	}

	if err := config.WriteEmpty(path, config.Filename); err != nil {
		return err
	}

	return nil
}

// ExampleProject creates the verless example project in the
// given path, containing all files from the example directory.
func ExampleProject(path string) error {
	for file, content := range files {
		dir := filepath.Dir(file)

		if err := fs.MkdirAll(path, dir); err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath.Join(path, file), []byte(content), 0755); err != nil {
			return err
		}
	}

	return nil
}
