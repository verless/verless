package atom

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"github.com/verless/verless/model"
)

const (
	filename string = "atom.xml"
)

func New(meta *model.Meta, outputDir string) *atom {
	a := atom{
		meta: meta,
		feed: &feeds.Feed{
			Title:       meta.Title,
			Link:        &feeds.Link{Href: meta.Base},
			Description: meta.Description,
			Author:      &feeds.Author{Name: meta.Author},
			Updated:     time.Time{},
			Created:     time.Now(),
			Subtitle:    meta.Subtitle,
		},
		outputDir: outputDir,
	}

	return &a
}

type atom struct {
	meta      *model.Meta
	feed      *feeds.Feed
	outputDir string
}

func (a *atom) ProcessPage(route string, page *model.Page) error {
	if page.Hide {
		return nil
	}

	canonical := fmt.Sprintf("%s%s%s", a.meta.Base, route, page.ID)

	item := &feeds.Item{
		Title:       page.Title,
		Link:        &feeds.Link{Href: canonical},
		Description: page.Description,
		Id:          canonical,
		Created:     page.Date,
	}

	a.feed.Add(item)
	return nil
}

func (a *atom) Finalize() error {
	path := filepath.Join(a.outputDir, filename)
	atomFile, err := os.Create(path)
	if err != nil {
		return err
	}

	return a.feed.WriteAtom(atomFile)
}
