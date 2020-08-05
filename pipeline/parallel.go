package pipeline

func runParallel(fn func(file string) error, files <-chan string, n int) error {
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
