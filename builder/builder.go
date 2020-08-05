package builder

import (
	"fmt"
	"github.com/verless/verless/model"
	"strings"
	"sync"
)

type walkFn func(route *model.Route) error

type Builder interface {
	RegisterPage(route string, page model.Page) error
	Dispatch() (model.Site, error)
}

func New() Builder {
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

	r := b.createRoute(route)
	r.Pages = append(r.Pages, page)

	return nil
}

func (b *builder) Dispatch() (model.Site, error) {
	return b.site, nil
}

func (b *builder) walkRoutes(walk walkFn, maxDepth int) error {
	return b.walkRoute(&b.site.Root, walk, maxDepth, 0)
}

func (b *builder) walkRoute(route *model.Route, walk walkFn, maxDepth, curDepth int) error {
	if maxDepth != -1 && curDepth == maxDepth {
		return nil
	}
	curDepth++

	if err := walk(route); err != nil {
		return err
	}

	for _, child := range route.Children {
		if err := b.walkRoute(child, walk, maxDepth, curDepth); err != nil {
			return err
		}
	}

	return nil
}

func (b *builder) createRoute(route string) *model.Route {
	var (
		node     = &b.site.Root
		segments = strings.Split(route, "/")
	)

	for i, s := range segments {
		if _, exists := node.Children[s]; !exists {
			node.Children[s] = &model.Route{
				Children:  make(map[string]*model.Route),
				Pages:     make([]model.Page, 0),
				IndexPage: model.IndexPage{},
			}
		}
		if i == len(segments)-1 {
			return node
		}
		node = node.Children[s]
	}

	return nil
}

func (b *builder) resolveRoute(route string) (*model.Route, error) {
	var (
		node     = &b.site.Root
		segments = strings.Split(route, "/")
	)

	for i, s := range segments {
		if i == len(segments)-1 {
			return node, nil
		}
		if _, exists := node.Children[s]; !exists {
			return node, fmt.Errorf("child route %s does not exist", s)
		}
		node = node.Children[s]
	}

	return node, fmt.Errorf("route %s does not exist", route)
}
