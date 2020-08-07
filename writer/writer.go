package writer

import (
	"github.com/verless/verless/model"
	"os"
	"path/filepath"
)

func New(path, outputDir string) *writer {
	w := writer{
		path:      path,
		outputDir: outputDir,
	}
	return &w
}

type writer struct {
	path      string
	outputDir string
}

func (w *writer) Write(site model.Site) error {
	_ = os.Mkdir(filepath.Join(w.path, w.outputDir), 0700)
	return nil
}
