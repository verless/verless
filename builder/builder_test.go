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

			segments := strings.Split(page.Route, "/")

			parent := builder.site.Root
			for i := 0; i < len(segments); i++ {
				test.NotEquals(t, nil, parent)

				if i == len(segments)-1 {
					test.Equals(t, page.ID, parent.Pages[0].ID)
					test.Assert(t, parent.IndexPage.Pages[0] == &parent.Pages[0], "the index page has to point to the actual page")
				} else {
					parent = *parent.Children[segments[i+1]]
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

	for name, _ := range tests {
		t.Log(name)
	}
}
