package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMessageAnyRouteError = "route '%v': %w"

	ErrWrongRouteFormat      = errors.New("the route has an invalid format")
	ErrChildNodeDoesNotExist = errors.New("child node does not exist")
)

// walkFn is invoked by WalkTree for each node in the route tree.
type walkFn func(node *Node) error

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

// WalkTree traverses the site's route tree and invokes the given
// walkFn on each node. Use maxDepth = -1 to traverse all nodes.
func (s *Site) WalkTree(walkFn walkFn, maxDepth int) error {
	return s.walkTreeNode(&s.Root, walkFn, maxDepth, 0)
}

// walkTreeNode invokes the walkFn on a given node and calls itself
// for all of its child nodes.
func (s *Site) walkTreeNode(node *Node, walkFn walkFn, maxDepth, curDepth int) error {
	if maxDepth != -1 && curDepth == maxDepth {
		return nil
	}
	curDepth++

	if err := walkFn(node); err != nil {
		return err
	}

	for _, child := range node.Children {
		if err := s.walkTreeNode(child, walkFn, maxDepth, curDepth); err != nil {
			return err
		}
	}

	return nil
}

// CreateNode creates a new node in the route tree. The route has
// to start with a slash representing the root route, e. g. /blog.
// Returns the error ErrMessageWrongRouteFormat if the given route has a invalid format.
// This is the case when the route does not start with a '/'.
func (s *Site) CreateNode(route string) (*Node, error) {
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
		s.Root.Children = make(map[string]*Node)
	}

	for i, s := range segments {
		if _, exists := node.Children[s]; !exists {
			node.Children[s] = &Node{
				Children: make(map[string]*Node),
				Pages:    make([]Page, 0),
				IndexPage: IndexPage{
					Page: Page{Route: route},
				},
			}
		}
		node = node.Children[s]
		if i == len(segments)-1 {
			return node, nil
		}
	}

	return nil, nil
}

// ResolveNode resolves and returns a route in the route tree.
// Has to start with a slash representing the root route.
// Returns the error ErrMessageWrongRouteFormat if the given route has a invalid format.
// This is the case when the route does not start with a '/'.
func (s *Site) ResolveNode(route string) (*Node, error) {
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
		s.Root.Children = make(map[string]*Node)
	}

	for i, s := range segments {
		if _, exists := node.Children[s]; !exists {
			return node, fmt.Errorf(ErrMessageAnyRouteError, s, ErrChildNodeDoesNotExist)
		}
		node = node.Children[s]
		if i == len(segments)-1 {
			return node, nil
		}
	}

	return node, fmt.Errorf("route %s does not exist", route)
}
