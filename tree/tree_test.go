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
	}

	for name, testCase := range tests {
		t.Log(name)

		root := &testNode{
			children: make(map[string]Node),
		}

		err := CreateNode(testCase.path, root, &testNode{})

		if testCase.isValidPath {
			test.Equals(t, nil, err)
		} else {
			test.NotEquals(t, nil, err)
		}

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
		path           string
		IsExistingPath bool
	}{
		"root path": {
			path:           RootPath,
			IsExistingPath: true,
		},
		"path with depth 1": {
			path:           "/blog",
			IsExistingPath: true,
		},
		"path with depth 2": {
			path:           "/blog/coffee",
			IsExistingPath: true,
		},
		"not existing path": {
			path: "/coffee",
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		node, err := ResolveNode(testCase.path, &root)

		if testCase.IsExistingPath {
			test.Equals(t, nil, err)
			test.NotEquals(t, nil, node)
		} else {
			test.NotEquals(t, nil, err)
		}
	}
}
