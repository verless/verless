package builder

import (
	"github.com/verless/verless/model"
	"sync"
)

func New() *builder {
	b := builder{
		site:  model.Site{},
		mutex: &sync.Mutex{},
	}
	return &b
}

type builder struct {
	site  model.Site
	mutex *sync.Mutex
}

func (b *builder) RegisterPage(route string, page model.Page) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	r := b.site.CreateRoute(route)
	r.Pages = append(r.Pages, page)

	return nil
}

func (b *builder) Dispatch() (model.Site, error) {
	return b.site, nil
}
