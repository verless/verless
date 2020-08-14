// Package fs provides functions for filesystem operations.
package fs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	// MarkdownOnly is a filter that only lets pass Markdown files.
	MarkdownOnly = func(file string) bool {
		return filepath.Ext(file) == ".md"
	}

	// NoUnderscores is a predefined filter that doesn't let pass
	// files starting with an underscore.
	NoUnderscores = func(file string) bool {
		filename := filepath.Base(file)
		return !strings.HasPrefix(filename, "_")
	}
)

// StreamFiles sends files in a given path that match the given
// filters through the files channel. If a value from stopSignal
// is received, StreamFiles exits.
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
				// do nothing, just re-run the for again until an
				// of the cases passes or the file can be sent.
			}
		}

		return nil
	})

	close(files)
	return err
}
