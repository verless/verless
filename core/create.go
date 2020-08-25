package core

import (
	"fmt"
	"os"

	"github.com/verless/verless/core/create"
)

const (
	exampleDir = "my-coffee-blog"
)

// CreateProjectOptions represents options for creating
// a new verless project.
type CreateProjectOptions struct {
	Overwrite bool
}

// CreateExampleOptions represents options for creating
// the verless example project.
type CreateExampleOptions struct {
	Overwrite bool
}

// RunCreateProject creates a new verless project.
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

// RunCreateExample creates the verless example project inside
// the current directory.
func RunCreateExample(options CreateExampleOptions) error {
	if _, err := os.Stat(exampleDir); !os.IsNotExist(err) {
		if !options.Overwrite {
			return fmt.Errorf("%s already exists. use --overwrite to overwrite it", exampleDir)
		}
		if err := os.RemoveAll(exampleDir); err != nil {
			return err
		}
	}

	return create.ExampleProject(exampleDir)
}
