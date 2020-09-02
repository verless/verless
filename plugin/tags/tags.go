// Package tags provides and implements the tags plugin.
package tags

import (
	"path/filepath"

	"github.com/verless/verless/model"
)

const (
	// Key is the public plugin key.
	Key string = "tags"
	// tagsDir is the target directory for all tag directories.
	tagsDir string = "/tags"
)

// New creates a new tags plugin that uses templates from the given
// build path and outputs the tag directories to outputDir.
func New() *tags {
	t := tags{
		m: make(map[string]*model.IndexPage),
	}

	return &t
}

// tags is the actual tags plugin that maintains a map with all
// tags from all processed pages.
type tags struct {
	m map[string]*model.IndexPage
}

// ProcessPage creates a new map entry for each tag in the processed
// page and adds the page to the entry's index page.
func (t *tags) ProcessPage(page *model.Page) error {
	for _, tag := range page.Tags {
		if _, exists := t.m[tag]; !exists {
			t.createIndexPage(tag)
		}
		t.m[tag].Pages = append(t.m[tag].Pages, page)
	}

	return nil
}

// PreWrite registers each index page in the site model. Those index
// pages will be rendered by the writer.
func (t *tags) PreWrite(site *model.Site) error {
	_, err := site.CreateNode(tagsDir)
	if err != nil {
		return err
	}

	for tag, indexPage := range t.m {
		node, err := site.CreateNode(filepath.Join(tagsDir, tag))
		if err != nil {
			return err
		}
		node.IndexPage = *indexPage
	}

	return nil
}

// PostWrite isn't needed by the tags plugin.
func (t *tags) PostWrite() error {
	return nil
}

// createIndexPage initializes a new index page for a given key.
func (t *tags) createIndexPage(key string) {
	t.m[key] = &model.IndexPage{
		Pages: make([]*model.Page, 0),
	}
}
