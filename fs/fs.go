package fs

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	MarkdownOnly = func(file string) bool {
		return filepath.Ext(file) == ".md"
	}
	NoUnderscores = func(file string) bool {
		filename := filepath.Base(file)
		return !strings.HasPrefix(filename, "_")
	}
)

func StreamFiles(path string, files chan<- string, filters ...func(file string) bool) error {
	err := filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, filter := range filters {
			if !filter(file) {
				return nil
			}
		}
		files <- file
		return nil
	})

	close(files)
	return err
}
