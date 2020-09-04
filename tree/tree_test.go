package tree

import "testing"

type testNode struct {
	children map[string]Node
}

func (t *testNode) Children() map[string]Node {
	return t.children
}

func (t *testNode) CreateChild(edge string) {
	child := &testNode{
		children: make(map[string]Node),
	}
	t.children[edge] = child
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
	root := &testNode{
		children: make(map[string]Node),
	}

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

		_ = CreateNode(testCase.path, root, &testNode{})
	}
}
