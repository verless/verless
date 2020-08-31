package builder

import (
	"strings"
	"testing"

	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
	"github.com/verless/verless/test"
)

var (
	// testPages is a set of pages used for testing.
	testPages []model.Page = []model.Page{
		{ID: "page-0", Route: "/route-0"},
		{ID: "page-1", Route: "/route-1"},
		{ID: "page-2", Route: "/route-2"},
		{ID: "page-3", Route: "/route-3"},
	}
)

// TestBuilder_RegisterPage checks if the testPages can be resolved
// from the site model exactly like they've been registered.
func TestBuilder_RegisterPage(t *testing.T) {
	tests := map[string]struct {
		pages []model.Page
	}{
		"register one testPage": {
			pages: testPages[:1],
		},
		"register several testPages": {
			pages: testPages,
		},
		"register with child pages": {
			pages: append(testPages,
				model.Page{ID: "child-0", Route: "/route-0/child-0"},
				model.Page{ID: "child-1", Route: "/route-0/child-0/child-1"},
				model.Page{ID: "child-2", Route: "/route-1/child-2"},
			),
		},
		"register with child pages whose parent doesn't exist": {
			pages: append(testPages,
				model.Page{ID: "child-1", Route: "/route-0/child-0/child-1"},
			),
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		builder := New(&config.Config{})

		for _, page := range testCase.pages {
			if err := builder.RegisterPage(page); err != nil {
				t.Fatal(err)
			}
		}

		for _, page := range testCase.pages {
			t.Logf("route %v", page.Route)

			routes := strings.Split(page.Route, "/")[1:]

			parent := builder.site.Root.Children[routes[0]]
			for i := 0; i < len(routes); i++ {
				test.NotEquals(t, nil, parent)

				if i == len(routes)-1 {
					test.Equals(t, page.ID, parent.Pages[0].ID)
					test.Assert(t, parent.IndexPage.Pages[0] == &parent.Pages[0], "the index page has to point to the actual page")
				} else {
					parent = parent.Children[routes[i+1]]
				}
			}
		}
	}
}

// TestBuilder_Dispatch checks if the dispatched site model is
// valid and contains all registered testPages.
func TestBuilder_Dispatch(t *testing.T) {
	tests := map[string]struct {
	}{
		// no tests as there is no logic yet
	}

	site, err := b.Dispatch()
	if err != nil {
		t.Fatal(err)
	}

	for i, page := range pages {
		segment := strings.TrimLeft(getRoute(i), "/")

		if site.Root.Children == nil {
			t.Fatalf("root route has uninitialized children map")
		}
		if _, exists := site.Root.Children[segment]; !exists {
			t.Fatalf("child route %s does not exist", segment)
		}

		route := site.Root.Children[segment]

		if len(route.Pages) < 1 {
			t.Fatalf("route %s contains no pages", segment)
		}
		if route.Pages[0].ID != page.ID {
			t.Errorf("expected page %s in route %s, got %s",
				page.ID, segment, route.Pages[0].ID)
		}
	}
}

	for name, _ := range tests {
		t.Log(name)
	}
}

// getRoute returns a generated route identified by a number n.
func getRoute(n int) string {
	return fmt.Sprintf("/route-%v", n)
}
