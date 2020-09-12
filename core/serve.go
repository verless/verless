package core

import (
	"log"
	"net"
	"sync"

	"github.com/spf13/afero"
	"github.com/verless/verless/config"
	"github.com/verless/verless/core/serve"
	"github.com/verless/verless/core/watch"
)

// ServeOptions represents options for running a verless serve command.
type ServeOptions struct {

	// BuildOptions stores all options for re-builds when watching the site.
	BuildOptions

	// Port specifies the port to run the server at.
	Port uint16

	// IP specifies the IP to listen on in combination with the port.
	IP net.IP

	// Watch enables automatic re-builds when a file changes.
	Watch bool
}

// RunServe serves a verless project using a simple file server.
// It can build the project automatically if ServeOptions.Build is true and
// even watch the whole project directory for changes if ServeOptions.Watch is true.
func RunServe(path string, options ServeOptions) error {
	// First check if the passed path is a verless project (valid verless cfg).
	_, err := config.FromFile(path, config.Filename)
	if err != nil {
		return err
	}

	targetFiles := finalOutputDir(path, &options.BuildOptions)

	// If yes, build it if requested to do so.
	options.BuildOptions.RecompileTemplates = options.Watch

	memMapFs := afero.NewMemMapFs()

	done := make(chan bool)
	if options.Watch {
		rebuildCh := make(chan string)

		// Only watch if needed.
		if options.Watch {
			if err := watch.Run(watch.Context{
				IgnorePath: targetFiles,
				Path:       path,
				ChangedCh:  rebuildCh,
				StopCh:     done,
			}); err != nil {
				return err
			}
		}

		var initialBuild sync.WaitGroup

		// Start rebuild goroutine.
		// If watch is not enabled, it's still used for the initial build.
		go func() {
			first := true
			initialBuild.Add(1)

			for {
				select {
				case _, ok := <-rebuildCh:
					if !ok {
						return
					}
					log.Println("rebuild")
					// Re-read config as it may have changed also.
					cfg, err := config.FromFile(path, config.Filename)
					if err != nil {
						log.Println("rebuild error:", err)
						continue
					}
					err = RunBuild(memMapFs, path, options.BuildOptions, cfg)
					if err != nil {
						log.Println("rebuild error:", err)
					}

					if first {
						initialBuild.Done()
						first = false
					}
				case _, ok := <-done:
					// Stops the goroutine if requested to.
					// Triggers on closing of the done channel.
					if !ok {
						return
					}
				}
			}
		}()

		// Trigger and wait for the initial rebuild.
		rebuildCh <- path
		initialBuild.Wait()

		// Stop rebuilding goroutine if not watching.
		if !options.Watch {
			done <- true
			close(done)
		}
	}

	// If the target folder doesn't exist, return an error.
	if _, err := memMapFs.Stat(targetFiles); err != nil {
		return err
	}

	// Then serve it.
	err = serve.Run(memMapFs, serve.Context{Path: targetFiles, Port: options.Port, IP: options.IP})

	// Stop building goroutine just to be sure.
	done <- true
	close(done)

	return err
}
