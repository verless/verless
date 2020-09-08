package model

import (
	"github.com/verless/verless/tree"
)

// Node represents an URL like /blog that contains multiple pages,
// an overview page (IndexPage) and child routes.
type Node struct {
	children  map[string]tree.Node
	Pages     []Page
	IndexPage IndexPage
}

// Children returns all children of a node.
func (n *Node) Children() map[string]tree.Node {
	return n.children
}

// InitChild initializes an empty node and links it to the current
// node by the given edge name.
func (n *Node) InitChild(edge string) {
	n.children[edge] = NewNode()
}

// CreateChild links a child node to the current node by the given
// edge name. If there's already a node behind the edge name, it will
// be overwritten.
func (n *Node) CreateChild(edge string, child tree.Node) {
	n.children[edge] = child
}

func NewNode() *Node {
	node := Node{
		children: map[string]tree.Node{},
	}
	return &node
}
