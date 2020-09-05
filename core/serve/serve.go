// Package serve provides verless' core serve functionality.
package serve

// Context provides all components required for serving an already built project.
type Context struct {
	Path string
	Port uint16
}

func Run(ctx Context) []error {
	// First check if the passed path is a verless project.
	// If yes, check if it has been already built (using the default target path)
	// If not build it.
	// Then serve it.

	return nil
}
