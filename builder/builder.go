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
		site: model.Site{
			Meta: cfg.Site.Meta,
		},
		mutex: &sync.Mutex{},
	}
	return &b
}

// builder represents a model builder maintaining a site model
// where all pages get registered.
type builder struct {
	site  model.Site
	mutex *sync.Mutex
}

// RegisterPage registers the given page under the given route. It
// is safe for concurrent usage.
func (b *builder) RegisterPage(route string, page model.Page) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var r *model.Route

	if route != "" {
		r = b.site.CreateRoute(route)
	} else {
		r = &b.site.Root
	}

	r.Pages = append(r.Pages, page)
	// ToDo: Append page to route's IndexPage

	return nil
}

// Dispatch finishes the model build and returns the model.
func (b *builder) Dispatch() (model.Site, error) {
	return b.site, nil
}
