package tree

import (
	"errors"
	"testing"

	"github.com/verless/verless/test"
)

type testNode struct {
	children map[string]Node
}

// Children simply returns the entire children map.
func (t *testNode) Children() map[string]Node {
	return t.children
}

// InitChild initializes a new child node on its own and registers
// it as a child linked to the current node by the given edge.
func (t *testNode) InitChild(edge string) {
	child := &testNode{
		children: make(map[string]Node),
	}
	t.children[edge] = child
}

// CreateChild registers a node as a child linked to the current
// node by the given edge.
func (t *testNode) CreateChild(edge string, node Node) {
	t.children[edge] = node
}

var root = testNode{
	children: map[string]Node{
		"blog": &testNode{
			children: map[string]Node{
				"coffee": &testNode{},
			},
		},
		"tags": &testNode{
			children: map[string]Node{
				"coffee": &testNode{},
			},
		},
	},
}

// TestCreateNode checks if the CreateNode registers all nodes
// under the given tree paths correctly.
func TestCreateNode(t *testing.T) {
	tests := map[string]struct {
		path string
	}{
		"root path": {
			path: RootPath,
		},
		"path with depth 1": {
			path: "/blog",
		},
		"path with depth 2": {
			path: "/blog/coffee",
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		root := &testNode{
			children: make(map[string]Node),
		}

		err := CreateNode(testCase.path, root, &testNode{})
		test.Ok(t, err)

		var n Node = root

		for _, edge := range Edges(testCase.path) {
			_, exists := n.Children()[edge]
			test.Equals(t, true, exists)
			n = n.Children()[edge]
		}
	}
}

// TestResolveOrInitNode tests if the ResolveOrInitNode function
// correctly initializes all nodes in a given path.
func TestResolveOrInitNode(t *testing.T) {
	tests := map[string]struct {
		path string
	}{
		"existing path with depth 1": {
			path: "/blog",
		},
		"existing path with depth 2": {
			path: "/blog/coffee",
		},
		"not existing path": {
			path: "/blog/espresso",
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		root := &testNode{
			children: make(map[string]Node),
		}

		node, err := ResolveOrInitNode(testCase.path, root)
		test.Ok(t, err)
		test.NotEquals(t, nil, node)

		var n Node = root

		for _, edge := range Edges(testCase.path) {
			_, exists := n.Children()[edge]
			test.Equals(t, true, exists)
			n = n.Children()[edge]
		}
	}
}

// TestResolveNode checks if the ResolveNode function can resolve
// nodes in a pre-defined tree correctly and returns the expected
// errors for non-existing paths.
func TestResolveNode(t *testing.T) {
	tests := map[string]struct {
		path          string
		expectedError error
	}{
		"root path": {
			path: RootPath,
		},
		"path with depth 1": {
			path: "/blog",
		},
		"path with depth 2": {
			path: "/blog/coffee",
		},
		"not existing path": {
			path:          "/coffee",
			expectedError: ErrEdgeNotFound,
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		_, err := ResolveNode(testCase.path, &root)
		if test.ExpectedError(t, testCase.expectedError, err) != test.IsCorrectNil {
			continue
		}
	}
}

// TestWalk checks if the Walk function invokes the walkFn for as
// many nodes as expected, depending on a given maxWidth.
func TestWalk(t *testing.T) {
	testErr := errors.New("this is a test error")

	tests := map[string]struct {
		depth         int
		count         int
		expectedError error
	}{
		"nodes with max depth 1": {
			depth: 1,
			count: 3,
		},
		"nodes with max depth 2": {
			depth: 2,
			count: 5,
		},
		"all nodes": {
			depth: -1,
			count: 5,
		},
		"with error": {
			depth:         -1,
			count:         0,
			expectedError: testErr,
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		count := 0

		err := Walk(&root, func(_ string, node Node) error {
			count++
			return testCase.expectedError
		}, testCase.depth)

		if test.ExpectedError(t, testCase.expectedError, err) != test.IsCorrectNil {
			continue
		}

		test.Equals(t, testCase.count, count)
	}
}

// TestWalkPath tests if the WalkPath correctly invokes the
// walkFn for each node in a given path.
func TestWalkPath(t *testing.T) {
	testErr := errors.New("this is a test error")

	tests := map[string]struct {
		path          string
		count         int
		expectedError error
	}{
		"root path": {
			path:  RootPath,
			count: 1,
		},
		"path with depth 1": {
			path:  "/blog",
			count: 2,
		},
		"path with depth 2": {
			path:  "/blog/coffee",
			count: 3,
		},
		"not existing path": {
			path:          "/blog/espresso",
			expectedError: testErr,
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		count := 0

		err := WalkPath(testCase.path, &root, func(node Node) error {
			count++
			return testCase.expectedError
		})

		if test.ExpectedError(t, testCase.expectedError, err) != test.IsCorrectNil {
			continue
		}

		test.Equals(t, testCase.count, count)
	}
}
