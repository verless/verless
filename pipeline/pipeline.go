package pipeline

import (
	"github.com/verless/verless/builder"
	"github.com/verless/verless/parser"
	"github.com/verless/verless/plugin"
	"github.com/verless/verless/writer"
)

type Pipeline interface {
	SetFileChan(files <-chan string)
	SetParser(parser parser.Parser)
	SetBuilder(builder builder.Builder)
	AddPlugin(plugin plugin.Plugin)
	SetWriter(writer writer.Writer)
	Run() error
}
