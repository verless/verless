package build

import (
	"github.com/verless/verless/model"
)

const (
	parallelism = 4
)

type Parser interface {
	ParsePage(src []byte) (model.Page, error)
}

type Builder interface {
	RegisterPage(route string, page model.Page) error
	Dispatch() (model.Site, error)
}

type Plugin interface {
	ProcessPage(page *model.Page) error
	Finalize() error
}

type Writer interface {
	Write(site model.Site) error
}

type Context struct {
	files   <-chan string
	parser  Parser
	builder Builder
	plugins []Plugin
	writer  Writer
}

func Run(ctx Context) error {
	err := runParallel(func(file string) error {
		return nil
	}, ctx.files, parallelism)

	if err != nil {
		return err
	}

	site, err := ctx.builder.Dispatch()
	if err != nil {
		return err
	}

	for _, plugin := range ctx.plugins {
		if err := plugin.Finalize(); err != nil {
			return err
		}
	}

	if err := ctx.writer.Write(site); err != nil {
		return err
	}

	return nil
}
