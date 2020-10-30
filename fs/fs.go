// Package fs provides functions for filesystem operations.
package fs

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
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

	// ErrStreaming is returned from StreamFiles.
	ErrStreaming error = nil
)

// StreamFiles sends all relative file paths inside a given path that
// match the given filters through the files channel.
func StreamFiles(path string, files chan<- string, filters ...func(file string) bool) error {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			close(files)
			return nil
		}
		return err
	}

	// Convert to absolute path so that it does not make a difference if
	// the paths are in different formats. e.g. one "example/" and the
	// other one "./example"
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	ErrStreaming = filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
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

		// Convert to absolute path so that it does not make a difference if
		// the paths are in different formats. e.g. one "example/" and the
		// other one "./example"
		file, err = filepath.Abs(file)
		if err != nil {
			return err
		}

		files <- file[len(path):]

		return nil
	})

	close(files)
	return ErrStreaming
}

// Rmdir removes an entire directory along with its contents. If the
// directory does not exist, nothing happens.
func Rmdir(targetFs afero.Fs, path string) error {
	if _, err := targetFs.Stat(path); !os.IsNotExist(err) {
		return err
	}

	return targetFs.RemoveAll(path)
}

// CopyFromOS copies a given directory from the OS filesystem into
// another filesystem instance to the desired destination.
//
// If fileOnly is set to true, files will be copied directly into the
// destination directory without their directory structure inside src.
func CopyFromOS(targetFs afero.Fs, src, dest string, fileOnly bool) error {
	var (
		files   = make(chan string)
		errchan = make(chan error)
		err     error
	)

	go func() {
		err := StreamFiles(src, files)
		errchan <- err
	}()

	for file := range files {
		var destPath, srcPath string

		// ToDo: Check if joining the filepath is okay in terms of performance.
		srcPath = filepath.Join(src, file)

		if fileOnly {
			filename := filepath.Base(file)
			destPath = filepath.ToSlash(filepath.Join(dest, filename))
		} else {
			destPath = filepath.ToSlash(filepath.Join(dest, file))
		}

		if err := targetFs.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		srcFile, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		destFile, err := targetFs.Create(destPath)
		if err != nil {
			return err
		}

		if _, err = io.Copy(destFile, srcFile); err != nil {
			return err
		}

		_ = srcFile.Close()
		_ = destFile.Close()
	}

	err = <-errchan
	return err
}

// IsSafeToRemove determines if a directory can be removed safely.
func IsSafeToRemove(targetFs afero.Fs, path string, force bool) bool {
	if force {
		return true
	}
	_, err := targetFs.Stat(path)
	return os.IsNotExist(err)
}
