package tags

import (
	"testing"

	"github.com/verless/verless/model"
	"github.com/verless/verless/test"
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
				test.NotEquals(t, nil, taggerTag)
				test.Assert(t, len(tagger.m[tag].Pages) > 0, "tag should exist")
			}
		}
	}
}

func TestTags_PreWrite(t *testing.T) {
	tests := map[string]struct {
		tagsIndexPages map[string]*model.IndexPage
		expectedError  error
	}{
		"normal index pages": {
			tagsIndexPages: map[string]*model.IndexPage{
				"test1": {},
				"test2": {},
				"test3": {},
			},
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		tagger := New()
		tagger.m = testCase.tagsIndexPages
		s := model.Site{}
		err := tagger.PreWrite(&s)
		if test.ExpectedError(t, testCase.expectedError, err) != test.IsCorrectNil {
			continue
		}

		tags, ok := s.Root.Children["tags"]
		test.Equals(t, true, ok)
		test.NotEquals(t, nil, tags)

		for tag := range tagger.m {
			child, ok := tags.Children[tag]
			test.Equals(t, true, ok)
			test.NotEquals(t, nil, child)
			test.NotEquals(t, nil, child.IndexPage)
			test.NotEquals(t, nil, child.IndexPage.Page)
		}
	}
}

func TestTags_PostWrite(t *testing.T) {}
