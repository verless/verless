package core

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
	"github.com/verless/verless/config"
	"github.com/verless/verless/out"
	"github.com/verless/verless/out/style"
	"github.com/verless/verless/theme"
)

// ServeOptions represents options for running a verless listenAndServe command.
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

// Serve serves a verless project using a simple file server.
// It can build the project automatically if ServeOptions.Build is true and
// even watch the whole project directory for changes if ServeOptions.Watch is true.
func Serve(path string, options ServeOptions) error {
	// First check if the passed path is a verless project (valid verless cfg).
	cfg, err := config.FromFile(path, config.Filename)
	if err != nil {
		return err
	}

	targetFiles := outputDir(path, &options.BuildOptions)

	// If yes, build it if requested to do so.
	options.BuildOptions.RecompileTemplates = options.Watch
	options.Overwrite = true

	memMapFs := afero.NewMemMapFs()

	done := make(chan bool)
	rebuildCh := make(chan string)

	// Only watch if needed.
	if options.Watch {
		if err := watch(watchContext{
			IgnorePaths: []string{
				targetFiles,
				filepath.Join(path, config.StaticDir, config.GeneratedDir),
				theme.GeneratedPath(path, cfg.Theme),
			},
			Path:      path,
			ChangedCh: rebuildCh,
			StopCh:    done,
		}); err != nil {
			return err
		}
	}

	var initialBuild sync.WaitGroup

	// Start rebuild goroutine.
	// If watch is not enabled, it's still used for the initial build.
	go func() {
		isFirst := true
		initialBuild.Add(1)

		for {
			select {
			case _, ok := <-rebuildCh:
				if !ok {
					return
				}
				out.T(style.Sparkles, "building project ...")

				build, err := NewBuild(memMapFs, path, options.BuildOptions)
				if err != nil {
					out.Err(style.Exclamation, "failed to initialize new build: %s", err.Error())
				}

				if err := build.Run(); err != nil {
					out.Err(style.Exclamation, "failed to build the project: %s", err.Error())
				}

				out.T(style.HeavyCheckMark, "project built successfully")

				if isFirst {
					initialBuild.Done()
					isFirst = false
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

	// If the target folder doesn't exist, return an error.
	if _, err := memMapFs.Stat(targetFiles); err != nil {
		return err
	}

	err = listenAndServe(memMapFs, targetFiles, options.IP, options.Port)

	// Stop building goroutine just to be sure.
	done <- true
	close(done)

	return err
}

// listenAndServe starts a file server serving the built project.
func listenAndServe(fs afero.Fs, path string, ip net.IP, port uint16) error {
	addr := fmt.Sprintf("%v:%v", ip, port)

	if ip.To4() == nil {
		addr = fmt.Sprintf("[%v]:%v", ip, port)
	}

	out.T(style.Bulb, "serving website on %s", addr)

	httpFs := afero.NewHttpFs(fs)
	server := http.FileServer(httpFs.Dir(path))
	http.Handle("/", server)

	return http.ListenAndServe(addr, server)
}
