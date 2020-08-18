package create

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
