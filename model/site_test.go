package model

import (
	"errors"
	"strings"
	"testing"

	"github.com/verless/verless/test"
)

var (
	// testPages is a set of pages used for testing.
	testPages []Page = []Page{
		{ID: "page-0", Route: "/route-0"},
		{ID: "page-1", Route: "/route-1"},
		{ID: "page-2", Route: "/route-2"},
		{ID: "page-3", Route: "/route-3"},
	}
)

// TestSite_CreateRoute checks if all routes are created correctly
// when pages are registered.
func TestSite_CreateRoute(t *testing.T) {
	tests := map[string]struct {
		route         string
		expectedError error
	}{
		"normal route without depth": {
			route: "/route-0",
		},
		"route not starting with /": {
			route:         "route-0",
			expectedError: ErrWrongRouteFormat,
		},
		"deeper route": {
			route: "/route-0/child-0/child1",
		},
		"empty route": {
			// Todo: currently this creates a child with key ""
			// is this expected?
			route: "/",
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		site := &Site{}

		actual, err := site.CreateRoute(testCase.route)
		if testCase.expectedError != nil {
			test.Assert(t, errors.Is(err, testCase.expectedError), "expected a specific error")
			continue
		}

		test.Ok(t, err)
		test.NotEquals(t, nil, actual)

		routes := strings.Split(testCase.route, "/")[1:]

		parent := &site.Root

		for i := 0; i <= len(routes); i++ {
			test.NotEquals(t, nil, parent)
			test.NotEquals(t, nil, parent.IndexPage)

			if i == len(routes) {
				test.Equals(t, 0, len(parent.Children))
			} else {
				test.NotEquals(t, 0, len(parent.Children))
				parent = parent.Children[routes[i]]
			}
		}
	}
}

/*

// TestSite_ResolveRoute checks if routes are resolvable from the
// route tree.
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

// TestSite_WalkRoutes checks if the walkFn is invoked for
// all nodes in the route tree.
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

// registerPages registers all pages in the site model.
func registerPages() {
	for i, page := range pages {
		route := s.CreateRoute(getRoute(i))
		route.Pages = append(route.Pages, page)
	}
}
*/
