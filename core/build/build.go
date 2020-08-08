package build

import (
	"fmt"
	"github.com/verless/verless/fs"
	"github.com/verless/verless/model"
	"io/ioutil"
)

const (
	parallelism = 4
)

type (
	Parser interface {
		ParsePage(src []byte) (model.Page, error)
	}
	Builder interface {
		RegisterPage(route string, page model.Page) error
		Dispatch() (model.Site, error)
	}
	Writer interface {
		Write(site model.Site) error
	}
	Plugin interface {
		ProcessPage(page *model.Page) error
		Finalize() error
	}
)

type Context struct {
	Path    string
	Parser  Parser
	Builder Builder
	Writer  Writer
	Plugins []Plugin
}

func Run(ctx Context) error {
	var (
		files  = make(chan string)
		errors = make(chan error)
	)

	go func() {
		errors <- fs.StreamFiles(ctx.Path, files, fs.MarkdownOnly, fs.NoUnderscores)
	}()

	err := runParallel(getFileProcessingPipeline(ctx), files, parallelism)

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

func getFileProcessingPipeline(ctx Context) func(file string) error {
	return func(file string) error {
		fmt.Println("read file", file)
		// read page
		fileData, err := ioutil.ReadFile(file)

		// parse page
		mdParsed, err := ctx.Parser.ParsePage(fileData)
		if err != nil {
			return err
		}

		fmt.Println("pared markdown:", file, mdParsed.Content)

		return nil
	}
}
