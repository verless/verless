package core

import (
	"fmt"
	"github.com/verless/verless/config"
	"os"
	"path/filepath"

	"github.com/verless/verless/core/create"
)

// CreateProjectOptions represents options for creating
// a new verless project.
type CreateProjectOptions struct {
	Overwrite bool
}

// RunCreateProject creates a new verless project. If the
// specified project path already exists, RunCreateProject
// returns an error unless --overwrite has been used.
func RunCreateProject(path string, options CreateProjectOptions) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if !options.Overwrite {
			return fmt.Errorf("%s already exists. use --overwrite to overwrite it", path)
		}
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return create.Project(path)
}

func RunCreateTheme(path, name string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist, use `verless create project NAME` first", path)
	}

	if _, err := os.Stat(filepath.Join(path, config.ThemesDir, name)); !os.IsNotExist(err) {
		return fmt.Errorf("a theme called %s already exists, remove it first", name)
	}

	return create.Theme(path, name)
}
