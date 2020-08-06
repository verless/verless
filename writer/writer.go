package writer

import (
	"github.com/verless/verless/model"
)

func New() *writer {
	w := writer{}
	return &w
}

type writer struct{}

func (w *writer) Write(site model.Site) error {
	return nil
}
