// Package builder provides a model builder implementation.
package builder

import (
	"sort"
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

	// If the page has been created as a file called index.md,
	// register the page as list page.
	if page.ID == config.ListPageID {
		node.ListPage.Page = page
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

	// Sort the pages of each node's list page by date.
	_ = tree.Walk(b.site.Root, func(node tree.Node) error {
		n := node.(*model.Node)

		sort.Slice(n.ListPage.Pages, func(i, j int) bool {
			return n.ListPage.Pages[i].Date.After(n.ListPage.Pages[j].Date)
		})

		return nil
	}, -1)

	return b.site, nil
}
