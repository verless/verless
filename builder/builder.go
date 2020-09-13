// Package builder provides a model builder implementation.
package builder

import (
	"sync"

	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
	"github.com/verless/verless/tree"
)

// New creates a new builder instance.
func New(cfg *config.Config) *builder {
	b := builder{
		site:  model.NewSite(),
		cfg:   cfg,
		mutex: &sync.Mutex{},
	}
	return &b
}

// builder represents a model builder maintaining a site model.
type builder struct {
	site  model.Site
	cfg   *config.Config
	mutex *sync.Mutex
}

// RegisterPage registers a given page under a given route. It
// is safe for concurrent usage.
func (b *builder) RegisterPage(page model.Page) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	n, err := tree.ResolveOrInitNode(page.Route, b.site.Root)
	if err != nil {
		return err
	}

	node := n.(*model.Node)

	// Register the page as list page if it was created as index.md.
	if page.ID == config.ListPageID {
		node.ListPage = model.ListPage{
			Page: page,
		}
	} else {
		// Otherwise, register the page as normal page.
		node.Pages = append(node.Pages, page)
		node.ListPage.Pages = append(node.ListPage.Pages, &node.Pages[len(node.Pages)-1])
	}

	if err := tree.CreateNode(page.Route, b.site.Root, node); err != nil {
		return err
	}

	return nil
}

// Dispatch finishes the model build and returns the model.
func (b *builder) Dispatch() (model.Site, error) {
	b.site.Meta = b.cfg.Site.Meta
	b.site.Nav = b.cfg.Site.Nav
	b.site.Footer = b.cfg.Site.Footer

	return b.site, nil
}
