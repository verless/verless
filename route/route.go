// Package route provides route-related convenience functions.
package route

import (
	"path/filepath"
	"strings"
)

const (
	delimiter string = "/"
	root      string = "/"
)

// FromPath converts a filepath to a valid verless route that can
// be registered in a route tree.
//
// A filepath like `content/blog/making-barista-quality-espresso.md`
// will be converted to `/blog/making-barista-quality-espresso`, for
// example.
func FromPath(contentDir, path string) string {
	route := filepath.ToSlash(filepath.Dir(path))
	route = route[len(contentDir):]

	return route
}

// ToSegments returns the individual segments for a route. Those
// segments represent the nodes for a route tree.
func ToSegments(route string) []string {
	if IsRoot(route) {
		return []string{}
	}

	segments := strings.Split(route, delimiter)

	return segments
}

// IsValid indicates whether a given route is valid or not.
func IsValid(route string) bool {
	return strings.HasPrefix(route, root)
}

// IsRoot indicates if a given route is the root route.
func IsRoot(route string) bool {
	return route == root
}
