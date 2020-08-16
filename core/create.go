package core

import (
	"fmt"
	"os"

	"github.com/verless/verless/core/create"
)

type CreateProjectOptions struct {
	Force bool
}

func RunCreateProject(path string, options CreateProjectOptions) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if !options.Force {
			return fmt.Errorf("%s already exists. use --force to overwrite it", path)
		}
	}

	return create.Project(path)
}
