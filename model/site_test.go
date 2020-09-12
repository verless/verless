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

	rootOnlyNode = Node{
		Children: map[string]*Node{},
	}
	complexNode = Node{
		Children: map[string]*Node{
			"child1": {
				Children: map[string]*Node{
					"child1-1": {
						Pages: []Page{
							{}, // any page
						},
						Children: map[string]*Node{
							"child1-2": {
								Pages: []Page{
									{}, // any page
								},
							},
						},
					},
				},
			},
			"child2": {
				Children: map[string]*Node{
					"child2-2": {
						Pages: []Page{
							{}, // any page
							{}, // any page
						},
					},
				},
				Pages: []Page{
					{}, // any page
				},
			},
			"child3": {},
		},
	}
)

// TestSite_CreateNode checks if all routes are created correctly
// when pages are registered.
func TestSite_CreateNode(t *testing.T) {
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
		"only root": {
			route: "/",
		},
		"completely empty route": {
			route:         "",
			expectedError: ErrWrongRouteFormat,
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		site := &Site{}

		actual, err := site.CreateNode(testCase.route)
		if test.ExpectedError(t, testCase.expectedError, err) != test.IsCorrectNil {
			continue
		}

		test.NotEquals(t, nil, actual)

		routes := strings.Split(testCase.route, "/")[1:]

		parent := &site.Root

		for i := 0; i <= len(routes); i++ {
			test.NotEquals(t, nil, parent)
			test.NotEquals(t, nil, parent.ListPage)

			// check also special case "" -> no child with "" should be created
			if i == len(routes) || routes[i] == "" {
				test.Equals(t, 0, len(parent.Children))
			} else {
				test.NotEquals(t, 0, len(parent.Children))
				parent = parent.Children[routes[i]]
			}
		}
	}
}

// TestSite_ResolveNode checks if routes are resolvable from the
// route tree.
func TestSite_ResolveNode(t *testing.T) {
	type routeTest struct {
		route         string
		expectedError error
	}

	tests := map[string]struct {
		route        Node
		routesToTest []routeTest
	}{
		"only root": {
			route: rootOnlyNode,
			routesToTest: []routeTest{
				{route: "/"},
				{route: "", expectedError: ErrWrongRouteFormat},
			},
		},
		"with children": {
			route: complexNode,
			routesToTest: []routeTest{
				{route: "", expectedError: ErrWrongRouteFormat},
				{route: "/"},
				{route: "/child1"},
				{route: "/child1/child1-1"},
				{route: "/child2"},
				{route: "/child2/child2-2"},
				{route: "/child3"},
			},
		},
		"wrong routes": {
			route: complexNode,
			routesToTest: []routeTest{
				{route: "/child5", expectedError: ErrChildNodeDoesNotExist},
				{route: "/child1/child5", expectedError: ErrChildNodeDoesNotExist},
				{route: "//", expectedError: ErrChildNodeDoesNotExist},
			},
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		site := &Site{
			Root: testCase.route,
		}

		for _, routeToTest := range testCase.routesToTest {
			t.Logf("\ttest route '%v'", routeToTest.route)
			route, err := site.ResolveNode(routeToTest.route)
			if test.ExpectedError(t, routeToTest.expectedError, err) == test.IsCorrectNil {
				test.NotEquals(t, nil, route)
			}
		}
	}
}

// TestSite_WalkTree checks if the walkFn is invoked for
// all nodes in the route tree, counts the found pages and checks if returned errors are handled.
func TestSite_WalkTree(t *testing.T) {
	aSimpleError := errors.New("an error")

	tests := map[string]struct {
		route             Node
		expectedError     error
		routeCount        int
		pageCount         int
		alternativeWalkFn func(node *Node) error
		depth             int
	}{
		"only root": {
			route:      rootOnlyNode,
			routeCount: 1,
			pageCount:  0,
			depth:      -1,
		},
		"with children": {
			route:      complexNode,
			routeCount: 7,
			pageCount:  5,
			depth:      -1,
		},
		"less depth": {
			route:      complexNode,
			routeCount: 4,
			pageCount:  1,
			depth:      2,
		},
		"throw error": {
			route: complexNode,
			alternativeWalkFn: func(node *Node) error {
				return aSimpleError
			},
			expectedError: aSimpleError,
			depth:         -1,
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		site := &Site{
			Root: testCase.route,
		}

		routeCount := 0
		pageCount := 0

		walkFn := testCase.alternativeWalkFn

		if walkFn == nil {
			walkFn = func(node *Node) error {
				routeCount++
				pageCount += len(node.Pages)
				return nil
			}
		}

		err := site.WalkTree(walkFn, testCase.depth)

		if test.ExpectedError(t, testCase.expectedError, err) != test.IsCorrectNil {
			continue
		}

		test.Equals(t, testCase.routeCount, routeCount)
		test.Equals(t, testCase.pageCount, pageCount)
	}
}
