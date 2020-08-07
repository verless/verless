package build

import (
	"github.com/verless/verless/fs"
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
	Path    string
	Parser  Parser
	Builder Builder
	Plugins []Plugin
	Writer  Writer
}

func Run(ctx Context) error {
	var (
		files  = make(chan string)
		errors = make(chan error)
	)

	go func() {
		errors <- fs.StreamFiles(ctx.Path, files, fs.MarkdownOnly, fs.NoUnderscores)
	}()

	err := runParallel(func(file string) error {
		return nil
	}, files, parallelism)

	if err != nil {
		return err
	}

	site, err := ctx.Builder.Dispatch()
	if err != nil {
		return err
	}

	if err := ctx.Writer.Write(site); err != nil {
		return err
	}

	for _, plugin := range ctx.Plugins {
		if err := plugin.Finalize(); err != nil {
			return err
		}
	}

	return nil
}
