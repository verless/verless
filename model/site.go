package model

import (
	"fmt"
	"strings"
)

type walkFn func(path string, route *Route) error

type Site struct {
	Meta   Meta
	Nav    Nav
	Root   Route
	Footer Footer
}

func (s *Site) WalkRoutes(walkFn walkFn, maxDepth int) error {
	return s.walkRoute("", &s.Root, walkFn, maxDepth, 0)
}

func (s *Site) walkRoute(path string, route *Route, walkFn walkFn, maxDepth, curDepth int) error {
	if maxDepth != -1 && curDepth == maxDepth {
		return nil
	}
	curDepth++

	if err := walkFn(path, route); err != nil {
		return err
	}

	for path, child := range route.Children {
		if err := s.walkRoute(path, child, walkFn, maxDepth, curDepth); err != nil {
			return err
		}
	}

	return nil
}

func (s *Site) CreateRoute(route string) *Route {
	var (
		node     = &s.Root
		segments = strings.Split(route, "/")
	)

	for i, s := range segments {
		if _, exists := node.Children[s]; !exists {
			node.Children[s] = &Route{
				Children:  make(map[string]*Route),
				Pages:     make([]Page, 0),
				IndexPage: IndexPage{},
			}
		}
		if i == len(segments)-1 {
			return node
		}
		node = node.Children[s]
	}

	return nil
}

func (s *Site) ResolveRoute(route string) (*Route, error) {
	var (
		node     = &s.Root
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
