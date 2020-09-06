// Package watch provides verless' ability to watch and rebuild a project
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
	DoneCh     <-chan bool
}

func Run(ctx Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {

			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if strings.HasPrefix(event.Name, ctx.IgnorePath) {
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					ctx.ChangedCh <- event.Name
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watch error:", err)
			case _, ok := <-ctx.DoneCh:
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
