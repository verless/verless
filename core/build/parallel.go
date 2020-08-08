package build

import "sync"

func runParallel(fn func(file string) error, files <-chan string, n int) error {
	errors := make(chan error)
	stopSignal := make(chan bool)

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			for file := range files {
				select {
				case stop, ok := <-stopSignal:
					if !ok || stop {
						// stopping is signalled
						return
					}
				default:
				}

				errors <- fn(file)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		if err != nil {
			// signal all still running go routines to stop also
			stopSignal <- true
			close(errors)
			return err
		}
	}

	return nil
}
