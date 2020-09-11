// Package serve provides verless' core serve functionality.
package serve

import (
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/afero"
)

// Context provides all components required for serving an already
// built project.
type Context struct {
	Path string
	Port uint16
	IP   net.IP
}

// Run starts a file server serving the already built project
// For this it uses the provided information of the context.
func Run(fs afero.Fs, ctx Context) error {
	addr := fmt.Sprintf("%v:%v", ctx.IP, ctx.Port)

	if ctx.IP.To4() == nil {
		addr = fmt.Sprintf("[%v]:%v", ctx.IP, ctx.Port)
	}

	httpFs := afero.NewHttpFs(fs)
	server := http.FileServer(httpFs.Dir(ctx.Path))
	http.Handle("/", server)

	return http.ListenAndServe(addr, server)
}
