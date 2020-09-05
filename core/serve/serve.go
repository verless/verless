// Package serve provides verless' core serve functionality.
package serve

import "fmt"

// Context provides all components required for serving an already built project.
type Context struct {
	Path string
	Port uint16
}

func Run(ctx Context) error {
	fmt.Println(ctx)
	return nil
}
