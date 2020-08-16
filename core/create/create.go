package create

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/verless/verless/config"
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
  # when re-building it.
  overwrite: true`
)

// Project creates a new verless project along with its mandatory
// directories and files.
func Project(path string) error {
	if err := createDirectories(path, config.ContentDir, config.TemplateDir, config.AssetDir); err != nil {
		return err
	}

	configFile := fmt.Sprintf("%s.yml", config.Filename)

	file, err := os.Create(filepath.Join(path, configFile))
	if err != nil {
		return err
	}

	if _, err := file.Write([]byte(defaultConfig)); err != nil {
		return err
	}

	return nil
}

// createDirectories creates all given directories inside path.
func createDirectories(path string, directories ...string) error {
	for _, directory := range directories {
		if err := os.MkdirAll(filepath.Join(path, directory), 0644); err != nil {
			return err
		}
	}
	return nil
}
