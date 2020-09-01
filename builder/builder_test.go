package builder

import (
	"fmt"
	"strings"
	"testing"

	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
)

var (
	// b is the builder instance used for testing.
	b *builder = nil
	// pages is a set of pages used for testing.
	pages []model.Page = []model.Page{
		{ID: "page-0"},
		{ID: "page-1"},
		{ID: "page-2"},
		{ID: "page-3"},
	}
)

// TestBuilder_RegisterPage checks if the pages can be resolved
// from the site model exactly like they've been registered.
func TestBuilder_RegisterPage(t *testing.T) {
	setupBuilder()

	for i, page := range pages {
		page.Route = getRoute(i)
		if err := b.RegisterPage(page); err != nil {
			t.Fatal(err)
		}
	}

	for i, page := range pages {
		node, err := b.site.ResolveNode(getRoute(i))
		if err != nil {
			t.Fatal(err)
		}
		if len(node.Pages) < 1 {
			t.Fatalf("route %s contains no pages", getRoute(i))
		}
		if node.Pages[0].ID != page.ID {
			t.Errorf("expected page %s in route %s, got %s",
				page.ID, getRoute(i), node.Pages[0].ID)
		}
	}
}

// TestBuilder_Dispatch checks if the dispatched site model is
// valid and contains all registered pages.
func TestBuilder_Dispatch(t *testing.T) {
	setupBuilder()

	for i, page := range pages {
		page.Route = getRoute(i)
		if err := b.RegisterPage(page); err != nil {
			t.Fatal(err)
		}
	}

	site, err := b.Dispatch()
	if err != nil {
		t.Fatal(err)
	}

	for i, page := range pages {
		// Skip the first page since its route is "/".
		if i == 0 {
			continue
		}

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

// setupBuilder initializes the builder if required.
func setupBuilder() {
	if b == nil {
		b = New(&config.Config{})
	}
}

// getRoute returns a generated route identified by a number n.
func getRoute(n int) string {
	if n == 0 {
		return "/"
	}
	return fmt.Sprintf("/route-%v", n)
}
