// Package builder provides a model builder implementation.
package builder

import (
	"path/filepath"
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
		cache: make(map[string]*model.Node),
	}
	return &b
}

// builder represents a model builder maintaining a site model.
type builder struct {
	site  model.Site
	cfg   *config.Config
	mutex *sync.Mutex
	cache map[string]*model.Node
}

// RegisterPage registers a given page under a given route. It
// is safe for concurrent usage.
func (b *builder) RegisterPage(page model.Page) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	node, err := b.nodeFromCache(page.Route)
	if err != nil {
		return err
	}

	// If the page has been created as a file called index.md,
	// register the page as list page.
	if page.ID == config.ListPageID {
		node.ListPage.Page = page
		return nil
	}

	// Otherwise, register the page as normal page.
	node.Pages = append(node.Pages, page)

	// Assign the list page route if it hasn't a route yet.
	if node.ListPage.Route == "" {
		node.ListPage.Route = page.Route
	}

	if node.ListPage.Href == "" {
		node.ListPage.Href = filepath.Join(page.Route, page.ID)
	}

	// Reference the new page in all parent nodes as well.
	err = tree.WalkPath(page.Route, b.site.Root, func(currentNode tree.Node) error {
		n := currentNode.(*model.Node)
		n.ListPage.Pages = append(n.ListPage.Pages, &node.Pages[len(node.Pages)-1])
		return nil
	})

	return err
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

// nodeFromCache loads a node from the cache. If the node isn't
// registered in the cache yet, nodeFromCache will load it from
// the route tree first.
//
// Normally, each page that gets registered would cause a node
// lookup in the tree. Caching the looked up nodes avoids this.
func (b *builder) nodeFromCache(path string) (*model.Node, error) {
	if _, exists := b.cache[path]; !exists {
		// If the node isn't in the cache yet, load it from the
		// tree, where it either gets resolved or initialized.
		n, err := tree.ResolveOrInitNode(path, b.site.Root)
		if err != nil {
			return nil, err
		}
		b.cache[path] = n.(*model.Node)
	}

	return b.cache[path], nil
}
