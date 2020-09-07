// Package tree provides functions for performing operations related
// to tree data structures.
//
// A website's page structure can be easily depicted as a tree, so these
// functions are used for storing pages in a tree structure as well as
// reading pages from it.
//
// The package exclusively uses terminology for tree data structures:
// https://en.wikipedia.org/wiki/Tree_(data_structure)#Terminology
package tree

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidPath indicates that a path is not formally valid.
	ErrInvalidPath = errors.New("path is not a valid tree path")
	// ErrEdgeNotFound indicates that an edge cannot be found.
	ErrEdgeNotFound = errors.New("edge does not exist")
)

// Node represents a tree node which contains a value and has zero or
// more child nodes connected via edges.
type Node interface {
	// Children should return a map of child nodes. The map keys are
	// expected to be edge names that link the node to its child.
	Children() map[string]Node

	// InitChild should initialize a new child node that is linked
	// to the current node via the given edge name.
	//
	// This implies that a node initialized with InitChild can be
	// loaded with Children afterwards, and is available under the
	// edge name passed to CreateChild.
	InitChild(edge string)

	// CreateChild should create a new child node that is linked to
	// the current node via the given edge name. Its behavior should
	// be analogous to InitChild.
	CreateChild(edge string, child Node)
}

// CreateNode stores a node under a given tree path. It follows the
// path starting from the root node and creates all required children
// using Node.CreateChild if they don't exist yet.
//
// Passing an invalid tree path to CreateNode will lead to undefined
// behavior. Check the path using IsValidPath first.
func CreateNode(path string, root Node, node Node) error {
	if !IsValidPath(path) {
		return fmt.Errorf("create node %s: %w", path, ErrInvalidPath)
	}
	if IsRootPath(path) {
		return nil
	}

	n := root
	edges := Edges(path)

	for i, edge := range edges {
		if _, exists := n.Children()[edge]; !exists {
			// If the current edge is the last one of the tree path, this
			// is the edge where the node has to be created.
			if i == len(edges)-1 {
				n.CreateChild(edge, node)
				return nil
			}
			n.InitChild(edge)
		}
		n = n.Children()[edge]
	}

	return nil
}

// ResolveNode follows the given path starting from the root node and
// traverses all child nodes until the last edge is reached.
//
// Returns the node linked to the last edge in the path or an error if
// it cannot be resolved. Passing an invalid tree path to ResolveNode
// will lead to undefined behavior. Check the path using IsValidPath
// first.
func ResolveNode(path string, root Node) (Node, error) {
	if !IsValidPath(path) {
		return nil, fmt.Errorf("resolve node %s: %w", path, ErrInvalidPath)
	}
	if IsRootPath(path) {
		return root, nil
	}

	n := root

	for _, edge := range Edges(path) {
		// Stop traversing the tree when an edge cannot be found.
		if _, exists := n.Children()[edge]; !exists {
			return nil, fmt.
				Errorf("resolve node %s: edge %s: %w", path, edge, ErrEdgeNotFound)
		}
		n = n.Children()[edge]
	}

	return n, nil
}

// Walk traverses all nodes in a tree, starting from the root route.
// For each node, walkFn will be invoked. As soon as an error arises
// in one of the walkFns, the error is handed up to the caller.
//
// maxDepth is counted starting from 0, which represents the root
// node. Set maxDepth to -1 to walk down the entire tree.
func Walk(root Node, walkFn func(node Node) error, maxDepth int) error {
	return walkNode(root, walkFn, maxDepth, 0)
}

func walkNode(node Node, walkFn func(node Node) error, maxDepth, curDepth int) error {
	if maxDepth != -1 && curDepth-1 == maxDepth {
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
