package model

import (
	"fmt"
	"strings"
)

// walkFn can be invoked for each node in the route tree.
type walkFn func(path string, node *Node) error

// Site represents the actual website. The site model is generated
// and populated with data and content during the website build.
//
// Any build.Writer implementation is capable of rendering this
// model as a static website.
type Site struct {
	Meta   Meta
	Nav    Nav
	Root   Node
	Footer Footer
}

// WalkRoutes traverses the site's route tree and invokes the given
// walkFn on each node. Use maxDepth = -1 to traverse all nodes.
func (s *Site) WalkRoutes(walkFn walkFn, maxDepth int) error {
	return s.walkRoute("", &s.Root, walkFn, maxDepth, 0)
}

// walkRoute invokes the walkFn on a given node and calls itself
// for all of its child nodes.
func (s *Site) walkRoute(path string, node *Node, walkFn walkFn, maxDepth, curDepth int) error {
	if maxDepth != -1 && curDepth == maxDepth {
		return nil
	}
	curDepth++

	if err := walkFn(path, node); err != nil {
		return err
	}

	for path, child := range node.Children {
		if err := s.walkRoute(path, child, walkFn, maxDepth, curDepth); err != nil {
			return err
		}
	}

	return nil
}

// CreateNode creates a new node in the route tree. The route has
// to start with a slash representing the root route, e. g. /blog.
func (s *Site) CreateNode(route string) *Node {
	if route == "/" {
		return &s.Root
	}

	var (
		node     = &s.Root
		segments = strings.Split(route[1:], "/")
	)

	if s.Root.Children == nil {
		s.Root.Children = make(map[string]*Node)
	}

	for i, s := range segments {
		if _, exists := node.Children[s]; !exists {
			node.Children[s] = &Node{
				Children:  make(map[string]*Node),
				Pages:     make([]Page, 0),
				IndexPage: IndexPage{},
			}
		}
		node = node.Children[s]
		if i == len(segments)-1 {
			return node
		}
	}

	return nil
}

// ResolveRoute resolves and returns a route in the route tree.
// Has to start with a slash representing the root route.
func (s *Site) ResolveRoute(route string) (*Node, error) {
	if route == "/" {
		return &s.Root, nil
	}

	var (
		node     = &s.Root
		segments = strings.Split(route[1:], "/")
	)

	if s.Root.Children == nil {
		s.Root.Children = make(map[string]*Node)
	}

	for i, s := range segments {
		if _, exists := node.Children[s]; !exists {
			return node, fmt.Errorf("child route %s does not exist", s)
		}
		node = node.Children[s]
		if i == len(segments)-1 {
			return node, nil
		}
	}

	return node, fmt.Errorf("route %s does not exist", route)
}
