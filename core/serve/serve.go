// Package serve provides verless' core serve functionality.
package serve

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// Context provides all components required for serving an already built project.
type Context struct {
	Path string
	Port uint16
	IP   net.IP
}

func Run(ctx Context) error {
	var uri string
	if ctx.IP.To4() == nil {
		uri = fmt.Sprintf("[%v]:%v", ctx.IP, ctx.Port)
	} else {
		uri = fmt.Sprintf("%v:%v", ctx.IP, ctx.Port)

	}

	fs := http.FileServer(http.Dir(ctx.Path))
	http.Handle("/", fs)

	log.Printf("verless serving '%v' on %v\n", ctx.Path, uri)
	err := http.ListenAndServe(uri, fs)
	return err
}
