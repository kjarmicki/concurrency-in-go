package main

func orDone(done, c <-chan any) chan any {
	outputStream := make(chan any)
	go func() {
		defer close(outputStream)
		for {
			select {
			case <-done:
				return
			case output, ok := <-c:
				if !ok {
					return
				}
				select {
				case outputStream <- output:
				case <-done:
				}
			}
		}
	}()
	return outputStream
}
