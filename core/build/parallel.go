package build

// processFiles takes a function for processing a file path. This
// function is called concurrently for each file in the channel.
//
// If one of the functions returns a non-nil error, the error is
// returned by processFiles.
func processFiles(fn func(file string) error, files <-chan string, n int) error {
	errors := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			for file := range files {
				errors <- fn(file)
			}
		}()
	}

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
