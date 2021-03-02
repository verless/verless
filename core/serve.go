package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

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

	out.T(style.Sparkles, "building project ...")

	build, err := NewBuild(memMapFs, path, options.BuildOptions)
	if err != nil {
		return err
	}

	if err := build.Run(); err != nil {
		return err
	}

	out.T(style.HeavyCheckMark, "project built successfully")

	// If --watch is enabled, launch a goroutine that handles rebuilds.
	if options.Watch {
		factory := func() (*Build, error) {
			return NewBuild(memMapFs, path, options.BuildOptions)
		}
		go watchAndRebuild(factory, rebuildCh, done)
	}

	// If the target folder doesn't exist, return an error.
	if _, err := memMapFs.Stat(targetFiles); err != nil {
		return err
	}

	err = listenAndServe(memMapFs, targetFiles, options.IP, options.Port)
	close(done)

	return err
}

// watchAndRebuild watches the project for changes and rebuilds the project
// once a change is detected. Any errors will be printed directly.
func watchAndRebuild(factory func() (*Build, error), rebuildCh <-chan string, doneCh <-chan bool) {
	for {
		select {
		case _, ok := <-rebuildCh:
			if !ok {
				return
			}
			out.T(style.Sparkles, "rebuilding project ...")

			build, err := factory()
			if err != nil {
				out.Err(style.Exclamation, "failed to initialize new build: %s", err.Error())
				continue
			}

			if err := build.Run(); err != nil {
				out.Err(style.Exclamation, "failed to build the project: %s", err.Error())
				continue
			}

			out.T(style.HeavyCheckMark, "project built successfully")
		case _, _ = <-doneCh:
			return
		}
	}
}

// listenAndServe starts a file server serving the built project.
func listenAndServe(fs afero.Fs, path string, ip net.IP, port uint16) error {
	addr := fmt.Sprintf("%v:%v", ip, port)

	if ip.To4() == nil {
		addr = fmt.Sprintf("[%v]:%v", ip, port)
	}

	httpFs := afero.NewHttpFs(fs)

	server := http.Server{
		Addr:    addr,
		Handler: http.FileServer(httpFs.Dir(path)),
	}

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt)

	out.T(style.Bulb, "serving website on %s", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			// In case the HTTP server cannot serve, just exit.
			out.Err(style.X, err.Error())
			os.Exit(1)
		}
	}()

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// Perform a graceful shutdown once the interrupt signal is received.
	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
