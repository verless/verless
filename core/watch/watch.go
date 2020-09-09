// Package watch provides verless' ability to watch a project
// and react to changes in a verless project.
package watch

import (
	"log"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

// Context provides all components required for serving an already built project.
type Context struct {
	Path       string
	IgnorePath string
	ChangedCh  chan<- string
	StopCh     <-chan bool
}

// Run watches a verless project for changes and writes the changed
// files to the passed Context.ChangedCh channel.
// To stop the watcher just close the Context.StopCh channel.
// Context.IgnorePath can be used to ignore a path inside the given Context.Path.
func Run(ctx Context) error {
	w := watcher.New()
	w.FilterOps(watcher.Write)

	//r := regexp.MustCompile("^\\.[^.]+$")
	//w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {

			select {
			case event, ok := <-w.Event:
				if !ok {
					return
				}

				// Avoid emitting an event if the (folder) ctx.IgnorePath itself gets created / removed (e.g. if the target folder gets deleted).
				if strings.HasPrefix(event.Path, ctx.IgnorePath) {
					continue
				}

				log.Println("STH HAPPEND!", event.Path)

				ctx.ChangedCh <- event.Path

			case err, ok := <-w.Error:
				if !ok {
					return
				}
				log.Println("watcher error:", err)

				//case _, ok := <-ctx.StopCh:
				//	// This case catches if the watching should be stopped.
				//	// It just watches for a closed channel.
				//	if !ok {
				//		return
				//	}
			}
		}
	}()

	if err := w.AddRecursive(ctx.Path); err != nil {
		return err
	}

	if err := w.Ignore(ctx.IgnorePath); err != nil {
		return err
	}

	go func() {
		log.Fatal(w.Start(time.Millisecond * 100))
	}()

	<-ctx.StopCh

	return nil
}
