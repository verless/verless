package create

import (
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
