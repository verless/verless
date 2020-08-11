package model

import (
	"fmt"
	"strings"
	"testing"
)

var (
	s     *Site  = nil
	pages []Page = []Page{
		{ID: "page-0"},
		{ID: "page-1"},
		{ID: "page-2"},
		{ID: "page-3"},
	}
)

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

func TestSite_ResolveRoute(t *testing.T) {
	setupSite()
	registerPages()

	for i := 0; i < len(pages); i++ {
		path := getRoute(i)

		route, err := s.ResolveRoute(path)
		if err != nil {
			t.Error(err)
		}

		if len(route.Pages) < 1 {
			t.Errorf("did not receive pages in route %s", path)
		}
	}
}

func TestSite_WalkRoutes(t *testing.T) {
	setupSite()
	registerPages()
	count := 0

	if err := s.WalkRoutes(func(path string, route *Route) error {
		if count != 0 && len(route.Pages) < 1 {
			return fmt.Errorf("did not receive pages in route %s", path)
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

func setupSite() {
	if s == nil {
		s = &Site{}
	}
}

func registerPages() {
	for i, page := range pages {
		route := s.CreateRoute(getRoute(i))
		route.Pages = append(route.Pages, page)
	}
}

func getRoute(n int) string {
	return fmt.Sprintf("/route-%v", n)
}
