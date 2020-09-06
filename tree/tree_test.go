package tree

import (
	"testing"

	"github.com/verless/verless/test"
)

type testNode struct {
	children map[string]Node
}

func (t *testNode) Children() map[string]Node {
	return t.children
}

func (t *testNode) InitChild(edge string) {
	child := &testNode{
		children: make(map[string]Node),
	}
	t.children[edge] = child
}

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
		edges := Edges(testCase.path)

		for _, edge := range edges {
			_, exists := n.Children()[edge]
			test.Equals(t, true, exists)
			n = n.Children()[edge]
		}
	}
}

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

func TestWalk(t *testing.T) {
	tests := map[string]struct {
		depth int
		count int
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
	}

	for name, testCase := range tests {
		t.Log(name)

		count := 0

		err := Walk(&root, func(node Node) error {
			count++
			return nil
		}, testCase.depth)

		test.Equals(t, nil, err)
		test.Equals(t, testCase.count, count)
	}
}
