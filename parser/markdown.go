package parser

import (
	"bytes"

	"github.com/verless/verless/model"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// NewMarkdown initializes and returns a new Markdown parser.
func NewMarkdown() *markdown {
	m := markdown{
		gm: goldmark.New(
			goldmark.WithExtensions(meta.Meta, highlighting.Highlighting),
		),
	}
	return &m
}

// markdown is an internal type that satisfies the build.Parser
// interface and thus can be used for retrieving model.Pages.
type markdown struct {
	gm goldmark.Markdown
}

// ParsePage converts the byte contents of a Markdown file to
// an instance of model.Page.
func (m *markdown) ParsePage(src []byte) (model.Page, error) {
	var (
		page model.Page
		buf  bytes.Buffer
		ctx  = parser.NewContext()
	)

	if err := m.gm.Convert(src, &buf, parser.WithContext(ctx)); err != nil {
		return page, err
	}

	page.Content = buf.String()
	page.Meta = make(map[string]string)
	metadata := meta.Get(ctx)

	readMetadata(metadata, &page)

	return page, nil
}
