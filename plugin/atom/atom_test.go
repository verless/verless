package atom

import (
	"fmt"
	"testing"

	"github.com/verless/verless/model"
)

var (
	a     *atom        = nil
	pages []model.Page = []model.Page{
		{ID: "page-0"},
		{ID: "page-1"},
		{ID: "page-2"},
		{ID: "page-3"},
	}
)

func TestAtom_ProcessPage(t *testing.T) {
	setupAtom()

	for i, page := range pages {
		if err := a.ProcessPage(getRoute(i), &page); err != nil {
			t.Fatal(err)
		}
	}

	if len(a.feed.Items) != len(pages) {
		t.Fatalf("expected %d stored pages, got %d", len(pages), len(a.feed.Items))
	}

	for i := 0; i < len(pages); i++ {
		item := a.feed.Items[i]

		if item.Title != pages[i].Title {
			t.Errorf("expected title %s, got %s", pages[i].Title, item.Title)
		}

		canonical := fmt.Sprintf("%s%s/%s", a.meta.Base, getRoute(i), pages[i].ID)

		if item.Link.Href != canonical {
			t.Errorf("expected link %s, got %s", canonical, item.Link.Href)
		}
	}
}

func setupAtom() {
	if a == nil {
		a = New(&model.Meta{
			Base: "https://example.com",
		}, "")
	}
}

func getRoute(n int) string {
	return fmt.Sprintf("/route-%v", n)
}
