package create

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/verless/verless/config"
	"github.com/verless/verless/fs"
)

// Project creates a new verless project.
func Project(path string) error {
	if err := fs.MkdirAll(path, config.AssetDir, config.ContentDir, config.TemplateDir); err != nil {
		return err
	}

	configFile := fmt.Sprintf("%s.yml", config.Filename)

	f, err := os.Create(filepath.Join(path, configFile))
	if err != nil {
		return err
	}

	_, _ = f.Write(nil)
	_ = f.Close()

	return nil
}
