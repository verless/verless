// Package tree provides functions for performing operations related
// to tree data structures. A website's page structure can be easily
// depicted as a tree, so these functions are used for storing pages
// in a tree structure as well as reading pages from it.
package tree

import "fmt"

// Node represents a tree node which contains a value. Each node has
// zero or more child nodes (children).
type Node interface {
	// Children should return a map of child nodes. The map keys are
	// expected to be edge names that link the node to its child.
	Children() map[string]Node

	// CreateChild should initialize a new child node that is linked
	// to the current node by the given edge name.
	//
	// This implies that a node initialized with CreateChild can be
	// loaded with Children afterwards, and is available under the edge
	// name passed to CreateChild:
	//
	//	n.CreateChild("c")
	//	c := n.Children()["c"]
	//
	// Handling edge naming collisions is up to the implementation.
	CreateChild(edge string)
}

// CreateNode stores a node under a given tree path. It follows the
// path starting from the root node and creates all required children
// using Node.CreateChild if they don't exist yet.
func CreateNode(path string, root Node, node Node) error {
	if !IsValidPath(path) {
		return fmt.Errorf("%s is not a valid tree path", path)
	}
	if IsRootPath(path) {
		return nil
	}

	n := root

	for _, edge := range Edges(path) {
		if _, exists := n.Children()[edge]; !exists {
			n.CreateChild(edge)
		}
		n = n.Children()[edge]
	}

	return nil
}

// ResolveNode follows the given path starting from the root node and
// traverses all child nodes until the last edge is reached.
//
// Returns the node linked to the last edge in the path or an error if
// it cannot be resolved.
func ResolveNode(path string, root Node) (Node, error) {
	if !IsValidPath(path) {
		return nil, fmt.Errorf("%s is not a valid tree path", path)
	}
	if IsRootPath(path) {
		return root, nil
	}

	n := root

	for _, edge := range Edges(path) {
		if _, exists := n.Children()[edge]; !exists {
			return nil, fmt.Errorf("edge %s does not exist", edge)
		}
		n = n.Children()[edge]
	}

	return n, nil
}

// Walk traverses all nodes in a tree, starting from the root route.
// For each node, walkFn will be invoked. As soon as an error arises
// in one of the walkFns, the error is handed up to the caller.
//
// If maxDepth is reached, Walk will return. If maxDepth is set to -1,
// Walk will walk down the entire tree.
func Walk(root Node, walkFn func(node Node) error, maxDepth int) error {
	return walkNode(root, walkFn, maxDepth, 0)
}

func walkNode(node Node, walkFn func(node Node) error, maxDepth, curDepth int) error {
	if curDepth == maxDepth {
		return nil
	}

	curDepth++

	if err := walkFn(node); err != nil {
		return err
	}

	for _, child := range node.Children() {
		if err := walkNode(child, walkFn, maxDepth, curDepth); err != nil {
			return err
		}
	}

	return nil
}
