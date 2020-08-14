package build

import "sync"

// processFiles takes a function for processing a file path. This
// function is called concurrently for each file in the channel.
//
// If a function returns an error it is passed to the errors channel
// and all goroutines stop accepting new files. But there may still
// be more than one error if two goroutines fail at the same time.
func processFiles(fn func(file string) error, files <-chan string, errors chan<- error, n int) {
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for file := range files {
				err := fn(file)
				if err != nil {
					errors <- fn(file)
				}

				// stop processing further files if any error occurred in any goroutine
				if len(errors) > 0 {
					break
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
