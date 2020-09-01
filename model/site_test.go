package model

import (
	"fmt"
	"strings"
	"testing"
)

var (
	// s is the Site instance used for testing.
	s *Site = nil
	// pages is a set of pages used for testing.
	pages []Page = []Page{
		{ID: "page-0"},
		{ID: "page-1"},
		{ID: "page-2"},
		{ID: "page-3"},
	}
)

// TestSite_CreateRoute checks if all routes are created correctly
// when pages are registered.
func TestSite_CreateRoute(t *testing.T) {
	setupSite()
	registerPages()

	for i := 0; i < len(pages); i++ {
		segment := strings.TrimLeft(getRoute(i), "/")

		if _, exists := s.Root.Children[segment]; !exists {
			t.Fatalf("route %s does not exist", segment)
		}

		route := s.Root.Children[segment]

		if len(route.Pages) < 1 {
			t.Errorf("no pages have been added to route %s", segment)
		}
	}
}

// TestSite_ResolveRoute checks if routes are resolvable from the
// route tree.
func TestSite_ResolveRoute(t *testing.T) {
	setupSite()
	registerPages()

	for i := 0; i < len(pages); i++ {
		route := getRoute(i)

		node, err := s.ResolveNode(route)
		if err != nil {
			t.Error(err)
		}

		if len(node.Pages) < 1 {
			t.Errorf("did not receive pages in route %s", route)
		}
	}
}

// TestSite_WalkTree checks if the walkFn is invoked for
// all nodes in the route tree.
func TestSite_WalkTree(t *testing.T) {
	setupSite()
	registerPages()
	count := 0

	if err := s.WalkTree(func(route string, node *Node) error {
		if count != 0 && len(node.Pages) < 1 {
			return fmt.Errorf("did not receive pages in node %s", route)
		}
		count++
		return nil
	}, -1); err != nil {
		t.Error(err)
	}

	if count-1 != len(pages) {
		t.Error("not all routes have been walked")
	}
}

// setupSite initializes the Site if required.
func setupSite() {
	if s == nil {
		s = &Site{}
	}
}

// registerPages registers all pages in the site model.
func registerPages() {
	for i, page := range pages {
		node := s.CreateNode(getRoute(i))
		node.Pages = append(node.Pages, page)
	}
}

// getRoute returns a generated route identified by a number n.
func getRoute(n int) string {
	return fmt.Sprintf("/route-%v", n)
}
