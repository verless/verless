package tags

import (
	"github.com/verless/verless/test"
	"testing"

	"github.com/verless/verless/model"
)

var (
	// testPages is a set of pages used for testing.
	testPages = []model.Page{
		{ID: "page-0", Route: "/route-0", Tags: []string{"t-1", "t-2"}},
		{ID: "page-1", Route: "/route-1", Tags: []string{"t-1", "t-3"}},
		{ID: "page-2", Route: "/route-2", Tags: []string{"t-2", "t-3"}},
		{ID: "page-3", Route: "/route-3", Tags: []string{"t-2"}},
	}
)

func TestTags_ProcessPage(t *testing.T) {
	tests := map[string]struct {
		pages         []model.Page
		expectedError error
	}{
		"normal pages": {
			pages: testPages,
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		tagger := New()

		for i, page := range testCase.pages {
			t.Logf("process page number %v, route '%v'", i, page.Route)
			err := tagger.ProcessPage(&page)
			if test.ExpectedError(t, testCase.expectedError, err) != test.IsCorrectNil {
				return
			}

			for _, tag := range page.Tags {
				taggerTag, exists := tagger.m[tag]

				test.Assert(t, exists, "tag should exist")
				test.Equals(t, nil, taggerTag)
				test.Assert(t, len(tagger.m[tag].Pages) > 0, "tag should exist")
			}
		}
	}
}

/*
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
			node, err := site.ResolveNode(filepath.Join(tagsDir, tag))
			if err != nil {
				t.Fatal(err)
			}
			if len(node.IndexPage.Pages) < 1 {
				t.Errorf("expected at least %d pages for tag %s, got %d", 1, tag, 0)
			}
		}
	}
}*/

func TestTags_PostWrite(t *testing.T) {}
