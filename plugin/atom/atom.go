// Package atom provides and implements the atom plugin.
package atom

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"github.com/spf13/afero"
	"github.com/verless/verless/model"
)

const (
	// filename is the filename for the RSS feed.
	filename string = "atom.xml"
)

// New creates a new atom plugin that generated a RSS feed with the
// provided metadata and stores the XML file in outputDir.
func New(meta *model.Meta, fs afero.Fs, outputDir string) *atom {
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
		fs:        fs,
		outputDir: outputDir,
	}

	a.feedItems = make(chan *feeds.Item)

	return &a
}

// atom is the actual atom plugin. It stores all RSS feed items
// as a feeds.Feed and renders those items in a XML file.s
type atom struct {
	meta      *model.Meta
	feed      *feeds.Feed
	feedItems chan *feeds.Item
	fs        afero.Fs
	outputDir string

	// workerShouldStop is a channel which indicates that the worker should stop working.
	// To stop it, pass true to it. The worker will close the channel.
	workerShouldStop chan bool
	// workerFinishedSignal gets closed by the worker to indicate that it finished all it's work.
	workerFinishedSignal chan bool
}

// PreProcessPages starts a worker goroutine which handles the a.feed.Add.
// This improves speed as the ProcessPage can add new items in a non blocking way.
func (a *atom) PreProcessPages() error {
	a.workerShouldStop = make(chan bool)
	a.workerFinishedSignal = make(chan bool)

	go func() {
		defer func() {
			close(a.workerShouldStop)
			close(a.workerFinishedSignal)
		}()

		for {
			select {
			case shouldStop := <-a.workerShouldStop:
				if shouldStop {
					return
				}
			default:
			}

			select {
			case feedItem := <-a.feedItems:
				a.feed.Add(feedItem)
			default:
			}
		}
	}()

	return nil
}

// ProcessPage takes a page to be processed by the plugin, reads
// metadata for that page and creates a new feed item from it.
func (a *atom) ProcessPage(page *model.Page) error {
	if page.Hidden || page.IsCustomListPage() {
		return nil
	}

	canonical := fmt.Sprintf("%s%s", a.meta.Base, page.Href)

	item := &feeds.Item{
		Title:       page.Title,
		Link:        &feeds.Link{Href: canonical},
		Description: page.Description,
		Id:          canonical,
		Created:     page.Date,
	}

	a.feedItems <- item
	return nil
}

// PostProcessPages stops and waits for the worker from PreProcessPages to finish.
func (a *atom) PostProcessPages() error {
	if a.workerShouldStop != nil {
		a.workerShouldStop <- true
	}
	_, _ = <-a.workerFinishedSignal
	return nil
}

// PreWrite isn't needed by the atom plugin.
func (a *atom) PreWrite(_ *model.Site) error {
	return nil
}

// PostWrite writes the internal feed.Feed instance into a file
// directly in the output directory.
func (a *atom) PostWrite() error {
	path := filepath.Join(a.outputDir, filename)
	atomFile, err := a.fs.Create(path)
	if err != nil {
		return err
	}

	return a.feed.WriteAtom(atomFile)
}
