package tree

import (
	"github.com/verless/verless/test"
	"testing"
)

// TestIsRootPath checks if the IsRootPath function correctly
// indicates whether a given path is a tree's root path or not.
func TestIsRootPath(t *testing.T) {
	tests := map[string]struct {
		path       string
		isRootPath bool
	}{
		"root path": {
			path:       RootPath,
			isRootPath: true,
		},
		"non-root path": {
			path: "/blog",
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		isRootPath := IsRootPath(testCase.path)
		test.Equals(t, testCase.isRootPath, isRootPath)
	}
}

// TestIsValidPath checks if the IsValidPath function correctly
// indicates whether a given path is valid or not.
func TestIsValidPath(t *testing.T) {
	tests := map[string]struct {
		path        string
		isValidPath bool
	}{
		"root path": {
			path:        RootPath,
			isValidPath: true,
		},
		"path with depth 1": {
			path:        "/blog",
			isValidPath: true,
		},
		"path with depth 2": {
			path:        "/blog/coffee",
			isValidPath: true,
		},
		"invalid path with depth 1": {
			path: "blog",
		},
		"invalid path with depth 2": {
			path: "blog/coffee",
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		isValidPath := IsValidPath(testCase.path)
		test.Equals(t, testCase.isValidPath, isValidPath)
	}
}

// TestEdges checks if the Edges function correctly returns all
// edges of a given tree path.
func TestEdges(t *testing.T) {
	tests := map[string]struct {
		path  string
		edges []string
	}{
		"root path": {
			path:  RootPath,
			edges: []string{},
		},
		"path with depth 1": {
			path:  "/blog",
			edges: []string{"blog"},
		},
		"path with depth 2": {
			path:  "/blog/coffee",
			edges: []string{"blog", "coffee"},
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		edges := Edges(testCase.path)
		test.Equals(t, testCase.edges, edges)
	}
}
