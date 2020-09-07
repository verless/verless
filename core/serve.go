package core

import (
	"log"
	"net"
	"os"

	"github.com/verless/verless/config"
	"github.com/verless/verless/core/serve"
	"github.com/verless/verless/core/watch"
)

// ServeOptions represents options for running a verless serve command.
type ServeOptions struct {
	BuildOptions
	// Port specifies the port to run the server at.
	Port uint16

	// IP specifies the IP to listen on in combination with the port.
	IP net.IP

	// Build enables automatic building of the verless project before serving.
	// Build is ignored when Watch is true.
	Build bool

	// Build enables automatic building of the verless project before and while serving.
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

	done := make(chan bool)
	if options.Build || options.Watch {
		rebuildCh := make(chan string)

		// Only watch if needed.
		if options.Watch {
			err := watch.Run(watch.Context{
				IgnorePath: targetFiles,
				Path:       path,
				ChangedCh:  rebuildCh,
				StopCh:     done,
			})

			if err != nil {
				return err
			}
		}

		firstBuildFinished := make(chan bool)

		// Start rebuild goroutine.
		// If watch is not enabled, it's still used for the initial build.
		go func() {
			first := true
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
					err = RunBuild(path, options.BuildOptions, cfg)
					if err != nil {
						log.Println("rebuild error:", err)
					}

					if first {
						firstBuildFinished <- true
						close(firstBuildFinished)
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

		// Trigger an initial rebuild.
		rebuildCh <- path

		// Stop rebuilding goroutine if not watching.
		if !options.Watch {
			done <- true
			close(done)
		}
		<-firstBuildFinished
	}

	// If the target folder doesn't exist, return an error.
	if _, err := os.Stat(targetFiles); err != nil {
		return err
	}

	// Then serve it.
	err = serve.Run(serve.Context{Path: targetFiles, Port: options.Port, IP: options.IP})

	// Stop building goroutine just to be sure.
	done <- true
	close(done)

	return err
}
