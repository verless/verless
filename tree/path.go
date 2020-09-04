package tree

import "strings"

const (
	// RootPath represents the path of the root node in a tree.
	RootPath  string = "/"
	delimiter string = "/"
)

// IsRootPath checks if a path corresponds to a tree's root path.
func IsRootPath(path string) bool {
	return path == RootPath
}

// IsValidPath checks if a path is considered valid by the tree
// package, meaning that it can be used in functions like CreateNode
// or ResolveNode.
func IsValidPath(path string) bool {
	return strings.HasPrefix(path, RootPath)
}

// Edges returns the tree edges for a path. Those edges will be
// used by the tree package to traverse a node's children.
func Edges(path string) []string {
	if IsRootPath(path) {
		return []string{}
	}

	return strings.Split(path, delimiter)
}
