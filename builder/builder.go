// Package builder provides a model builder implementation.
package builder

import (
	"sync"

	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
)

// New creates a new builder instance.
func New(cfg *config.Config) *builder {
	b := builder{
		site:  model.Site{},
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

	var r *model.Route

	if page.Route != "" {
		r = b.site.CreateRoute(page.Route)
	} else {
		r = &b.site.Root
	}

	r.Pages = append(r.Pages, page)
	r.IndexPage.Pages = append(r.IndexPage.Pages, &page)

	return nil
}

// Dispatch finishes the model build and returns the model.
func (b *builder) Dispatch() (model.Site, error) {
	b.site.Meta = b.cfg.Site.Meta
	b.addNavAndFooter()

	return b.site, nil
}

// addNavAndFooter creates the site's navigation and footer items.
func (b *builder) addNavAndFooter() {
	for _, i := range b.cfg.Site.Nav.Items {
		b.site.Nav.Items = append(b.site.Nav.Items, model.NavItem{
			Label:  i.Label,
			Target: i.Target,
		})
	}

	for _, i := range b.cfg.Site.Footer.Items {
		b.site.Footer.Items = append(b.site.Footer.Items, model.FooterItem{
			Label:  i.Label,
			Target: i.Target,
		})
	}
}
