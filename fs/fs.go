package fs

import (
	"errors"
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

func StreamFiles(path string, files chan<- string, stopSignal <-chan bool, filters ...func(file string) bool) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		close(files)
		return nil
	}

	err := filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		for _, filter := range filters {
			if !filter(file) {
				return nil
			}
		}

		for {
			if len(files) == cap(files) {
				files <- file
				break
			}

			select {
			case _, ok := <-stopSignal: // check for stop signal (channel closing)
				if !ok {
					return errors.New("forcefully stopped filepath walk")
				}
			default:
				// do nothing, just re-run the for again until any of the cases passes or the file can be sent.
			}
		}

		return nil
	})

	close(files)
	return err
}
