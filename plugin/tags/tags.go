// Package tags provides and implements the tags plugin.
package tags

import (
	"path/filepath"

	"github.com/verless/verless/model"
	"github.com/verless/verless/tree"
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
		m: make(map[string]*model.ListPage),
	}

	return &t
}

// tags is the actual tags plugin that maintains a map with all
// tags from all processed pages.
type tags struct {
	m map[string]*model.ListPage
}

// ProcessPage creates a new map entry for each tag in the processed
// page and adds the page to the entry's list page.
func (t *tags) ProcessPage(page *model.Page) error {
	for _, tag := range page.Tags {
		if _, exists := t.m[tag]; !exists {
			t.createListPage(tag)
		}
		t.m[tag].Pages = append(t.m[tag].Pages, page)
	}

	return nil
}

// PreWrite registers each list page in the site model. Those list
// pages will be rendered by the writer.
func (t *tags) PreWrite(site *model.Site) error {
	if err := tree.CreateNode(tagsDir, site.Root, model.NewNode()); err != nil {
		return err
	}

	for tag, listPage := range t.m {
		path := filepath.ToSlash(filepath.Join(tagsDir, tag))

		node := model.NewNode()
		node.ListPage = *listPage

		if err := tree.CreateNode(path, site.Root, node); err != nil {
			return err
		}
	}

	return nil
}

// PostWrite isn't needed by the tags plugin.
func (t *tags) PostWrite() error {
	return nil
}

// createListPage initializes a new list page for a given key.
func (t *tags) createListPage(key string) {
	t.m[key] = &model.ListPage{
		Pages: make([]*model.Page, 0),
	}
}
