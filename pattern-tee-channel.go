package main

import "fmt"

func PatternTeeChannel() {
	tee := func(done <-chan any, in <-chan any) (<-chan any, <-chan any) {
		out1 := make(chan any)
		out2 := make(chan any)

		go func() {
			defer close(out1)
			defer close(out2)

			for val := range orDone(done, in) {
				for i := 0; i < 2; i++ {
					var out1, out2 = out1, out2
					select {
					case <-done:
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
	}

	done := make(chan any)
	defer close(done)

	in := make(chan any)
	out1, out2 := tee(done, in)
	go func() {
		defer close(in)
		in <- 1
		in <- 2
		in <- 3
		in <- 4
	}()

	for val := range out1 {
		fmt.Printf("out1 %v, out2 %v\n", val, <-out2)
	}
}
