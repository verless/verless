// Package watch provides verless' ability to watch a project
// and react to changes in a verless project.
package watch

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

// Context provides all components required for serving an already built project.
type Context struct {
	Path        string
	IgnorePaths []string
	ChangedCh   chan<- string
	StopCh      <-chan bool
}

// Run watches a verless project for changes and writes the changed
// files to the passed Context.ChangedCh channel.
// To stop the watcher just close the Context.StopCh channel.
// Context.IgnorePath can be used to ignore a path inside the given Context.Path.
func Run(ctx Context) error {
	w := watcher.New()
	w.FilterOps(watcher.Write)

	go func() {
	watcherLoop:
		for {

			select {
			case event, ok := <-w.Event:
				if !ok {
					return
				}
				if filepath.Ext(event.Path) == "" {
					continue
				}

				for _, ignorePath := range ctx.IgnorePaths {
					p, err := filepath.Abs(ignorePath)
					if err != nil {
						log.Println(err)
						continue
					}
					if strings.HasPrefix(event.Path, p) {
						continue watcherLoop
					}
				}

				ctx.ChangedCh <- event.Path

			case err, ok := <-w.Error:
				if !ok {
					return
				}
				log.Println("watcher error:", err)

			case _, ok := <-ctx.StopCh:
				if !ok {
					w.Close()
				}
			}
		}
	}()

	if err := w.AddRecursive(ctx.Path); err != nil {
		return err
	}

	if err := w.Ignore(ctx.IgnorePaths...); err != nil {
		return err
	}

	var err error

	go func() {
		err = w.Start(time.Millisecond * 100)
	}()

	return err
}
