package tags

import (
	"fmt"
	"github.com/verless/verless/model"
	"path/filepath"
	"testing"
)

var (
	tg    *tags        = nil
	pages []model.Page = []model.Page{
		{ID: "page-0", Tags: []string{"t-1", "t-2"}},
		{ID: "page-1", Tags: []string{"t-1", "t-3"}},
		{ID: "page-2", Tags: []string{"t-2", "t-3"}},
		{ID: "page-3", Tags: []string{"t-2"}},
	}
)

func TestTags_ProcessPage(t *testing.T) {
	setupTags()

	for i, page := range pages {
		page.Route = getRoute(i)
		if err := tg.ProcessPage(&page); err != nil {
			t.Fatal(err)
		}
	}

	for _, page := range pages {
		for _, tag := range page.Tags {
			if _, exists := tg.m[tag]; !exists {
				t.Fatalf("key for tag %s does not exist", tag)
			}
			if tg.m[tag] == nil {
				t.Fatalf("IndexPage for tag %s is nil", tag)
			}
			if len(tg.m[tag].Pages) < 1 {
				t.Errorf("expected at least %d pages for tag %s, got %d", 1, tag, 0)
			}
		}
	}
}

func TestTags_PreWrite(t *testing.T) {
	setupTags()

	for i, page := range pages {
		page.Route = getRoute(i)
		if err := tg.ProcessPage(&page); err != nil {
			t.Fatal(err)
		}
	}

	site := model.Site{}

	if err := tg.PreWrite(&site); err != nil {
		t.Fatal(err)
	}

	for _, page := range pages {
		for _, tag := range page.Tags {
			route, err := site.ResolveRoute(filepath.Join(tagsDir, tag))
			if err != nil {
				t.Fatal(err)
			}
			if len(route.IndexPage.Pages) < 1 {
				t.Errorf("expected at least %d pages for tag %s, got %d", 1, tag, 0)
			}
		}
	}
}

func TestTags_PostWrite(t *testing.T) {}

func setupTags() {
	if tg == nil {
		tg = New("", "")
	}
}

func getRoute(n int) string {
	return fmt.Sprintf("/route-%v", n)
}
