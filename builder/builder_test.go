package builder

import (
	"fmt"
	"strings"
	"testing"

	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
)

var (
	b     *builder     = nil
	pages []model.Page = []model.Page{
		{ID: "page-0"},
		{ID: "page-1"},
		{ID: "page-2"},
		{ID: "page-3"},
	}
)

func TestBuilder_RegisterPage(t *testing.T) {
	setupBuilder()

	for i, page := range pages {
		if err := b.RegisterPage(getRoute(i), page); err != nil {
			t.Fatal(err)
		}
	}

	for i, page := range pages {
		route, err := b.site.ResolveRoute(getRoute(i))
		if err != nil {
			t.Fatal(err)
		}
		if len(route.Pages) < 1 {
			t.Fatalf("route %s contains no pages", getRoute(i))
		}
		if route.Pages[0].ID != page.ID {
			t.Errorf("expected page %s in route %s, got %s",
				page.ID, getRoute(i), route.Pages[0].ID)
		}
	}
}

func TestBuilder_Dispatch(t *testing.T) {
	setupBuilder()

	for i, page := range pages {
		if err := b.RegisterPage(getRoute(i), page); err != nil {
			t.Fatal(err)
		}
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

func setupBuilder() {
	if b == nil {
		b = New(&config.Config{})
	}
}

func getRoute(n int) string {
	return fmt.Sprintf("/route-%v", n)
}
