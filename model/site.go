package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMessageAnyRouteError = "route '%v': %w"

	ErrWrongRouteFormat       = errors.New("the route has an invalid format")
	ErrChildRouteDoesNotExist = errors.New("child route does not exist")
)

// walkFn can be invoked for each route in the route tree.
type walkFn func(path string, route *Route) error

// Site represents the actual website. The site model is generated
// and populated with data and content during the website build.
//
// Any build.Writer implementation is capable of rendering this
// model as a static website.
type Site struct {
	Meta   Meta
	Nav    Nav
	Root   Route
	Footer Footer
}

// WalkRoutes traverses the site's route tree and invokes the given
// walkFn on each node. Use maxDepth = -1 to traverse all nodes.
func (s *Site) WalkRoutes(walkFn walkFn, maxDepth int) error {
	return s.walkRoute("", &s.Root, walkFn, maxDepth, 0)
}

// walkRoute invokes the walkFn on a given route and calls itself
// for all of its child routes.
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

// CreateRoute creates a new route in the route tree. The route has to
// start with a slash representing the root route, e. g. /blog.
// Returns the error ErrMessageWrongRouteFormat if the given route has a invalid format.
func (s *Site) CreateRoute(route string) (*Route, error) {
	if !strings.HasPrefix(route, "/") {
		return nil, fmt.Errorf(ErrMessageAnyRouteError, route, ErrWrongRouteFormat)
	}

	if route == "/" {
		return &s.Root, nil
	}

	var (
		node     = &s.Root
		segments = strings.Split(route[1:], "/")
	)

	if s.Root.Children == nil {
		s.Root.Children = make(map[string]*Route)
	}

	for i, s := range segments {
		if _, exists := node.Children[s]; !exists {
			node.Children[s] = &Route{
				Children:  make(map[string]*Route),
				Pages:     make([]Page, 0),
				IndexPage: IndexPage{},
			}
		}
		node = node.Children[s]
		if i == len(segments)-1 {
			return node, nil
		}
	}

	return nil, nil
}

// ResolveRoute resolves and returns a route in the route tree.
// Has to start with a slash representing the root route.
func (s *Site) ResolveRoute(route string) (*Route, error) {
	if !strings.HasPrefix(route, "/") {
		return nil, fmt.Errorf(ErrMessageAnyRouteError, route, ErrWrongRouteFormat)
	}

	if route == "/" {
		return &s.Root, nil
	}

	var (
		node     = &s.Root
		segments = strings.Split(route[1:], "/")
	)

	if s.Root.Children == nil {
		s.Root.Children = make(map[string]*Route)
	}

	for i, s := range segments {
		if _, exists := node.Children[s]; !exists {
			return node, fmt.Errorf(ErrMessageAnyRouteError, s, ErrChildRouteDoesNotExist)
		}
		node = node.Children[s]
		if i == len(segments)-1 {
			return node, nil
		}
	}

	return node, fmt.Errorf("route %s does not exist", route)
}
