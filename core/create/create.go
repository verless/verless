package create

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// defaultConfig is the default YAML project configuration when
	// a new project is initialized.
	defaultConfig string = `# verless.yml is your project configuration.
# Check out the example project for a full configuration:
# https://github.com/verless/verless/tree/docs/example
site:
  # General information about your project.
  meta:
    title: 
    subtitle: 
    description: 
    author: 
    base: 
# Enable plugins for your project.
plugins:
  - atom
build:
  # Allow verless to overwrite the output directory of your website
  # when re-building it. False by default.
  overwrite: false`
)

// Project creates a new verless project which is an exact copy
// of the example project.
func Project(path string) error {
	for file, content := range files {
		dir := filepath.Join(path, filepath.Dir(file))

		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath.Join(path, file), []byte(content), 0755); err != nil {
			return err
		}
	}

	return nil
}
