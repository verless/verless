// Package watch provides verless' ability to watch a project
// and react to changes in a verless project.
package watch

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
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
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {

			select {
			case event, ok := <-watcher.Events:
				// This case catches if an event occurred to a watched file.
				if !ok {
					return
				}

				// Avoid emitting an event if the (folder) ctx.IgnorePath itself gets created / removed (e.g. if the target folder gets deleted).
				if strings.HasPrefix(event.Name, ctx.IgnorePath) {
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					ctx.ChangedCh <- event.Name
				}
			case err, ok := <-watcher.Errors:
				// This case catches if an error occurred while watching the files.
				if !ok {
					return
				}
				log.Println("watch error:", err)
			case _, ok := <-ctx.StopCh:
				// This case catches if the watching should be stopped.
				// It just watches for a closed channel.
				if !ok {
					return
				}
			}
		}
	}()

	err = filepath.Walk(ctx.Path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() || strings.HasPrefix(path, ctx.IgnorePath) {
			return err
		}

		err = watcher.Add(path)
		return err
	})

	return err
}
