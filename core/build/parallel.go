package build

import "sync"

func runParallel(fn func(file string) error, files <-chan string, n int) error {
	errors := make(chan error, n)
	stopSignals := make(chan bool, n)

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			for file := range files {
				select {
				case stop, ok := <-stopSignals:
					if !ok || stop {
						// stopping is signalled
						break
					}
				default:
				}

				// TODO: Also handle several errors in different files.
				// 		 Currently this leads to writing to a closed channel.
				//       Maybe pass an errors chanel from the outside?
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
			for i := 0; i < n; i++ {
				stopSignals <- true
			}

			close(errors)
			return err
		}
	}

	return nil
}
