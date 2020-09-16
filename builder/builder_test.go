package builder

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
	"github.com/verless/verless/test"
	"github.com/verless/verless/tree"
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

			segments := strings.Split(page.Route, "/")[1:]
			var parent tree.Node = builder.site.Root

			for i := 0; i < len(segments); i++ {
				test.NotEquals(t, nil, parent)
				parent = parent.Children()[segments[i]]
				if i == len(segments)-1 {
					test.Equals(t, page.ID, parent.(*model.Node).Pages[0].ID)
					test.Assert(t, parent.(*model.Node).ListPage.Pages[0] == &parent.(*model.Node).Pages[0],
						"the index page has to point to the actual page")
				}
			}
		}
	}
}

// TestBuilder_Dispatch checks if the Dispatch method returns the
// site model as expected.
func TestBuilder_Dispatch(t *testing.T) {
	tests := map[string]struct {
		site model.Site
	}{
		"site with multiple page references": {
			site: model.Site{
				Root: &model.Node{
					ListPage: model.ListPage{
						Pages: []*model.Page{
							{
								Date: time.Date(2020, 03, 16, 0, 0, 0, 0, time.UTC),
							},
							{
								Date: time.Date(2020, 01, 28, 0, 0, 0, 0, time.UTC),
							},
							{
								Date: time.Date(2020, 11, 3, 0, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		builder := New(&config.Config{})
		builder.site = testCase.site

		site, err := builder.Dispatch()
		test.Ok(t, err)

		// Test if all page references in list pages are sorted correctly.
		_ = tree.Walk(site.Root, func(node tree.Node) error {
			n := node.(*model.Node)

			isSorted := sort.SliceIsSorted(n.ListPage.Pages, func(i, j int) bool {
				return n.ListPage.Pages[i].Date.Before(n.ListPage.Pages[j].Date)
			})

			test.Assert(t, isSorted, "page references should be sorted")
			return nil
		}, -1)
	}
}
