package parser

import (
	"bytes"
	"github.com/verless/verless/model"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func NewMarkdown() *markdown {
	m := markdown{
		gm: goldmark.New(
			goldmark.WithExtensions(meta.Meta, highlighting.Highlighting),
		),
	}
	return &m
}

type markdown struct {
	gm goldmark.Markdown
}

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
	metadata := meta.Get(ctx)

	readMetadata(metadata, &page)

	return page, nil
}
