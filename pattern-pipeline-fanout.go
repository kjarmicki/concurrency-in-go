package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
)

func PatternPipelineFanout() {
	repeatFn := func(done <-chan any, fn func() any) <-chan any {
		outputStream := make(chan any)
		go func() {
			defer close(outputStream)
			for {
				select {
				case <-done:
					return
				case outputStream <- fn():
				}
			}
		}()
		return outputStream
	}

	toInt := func(done <-chan any, anyStream <-chan any) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for val := range anyStream {
				select {
				case <-done:
					return
				case intStream <- val.(int):
				}
			}
		}()
		return intStream
	}

	takeInt := func(done <-chan any, intStream <-chan int, limit int) <-chan int {
		outputStream := make(chan int)
		go func() {
			defer close(outputStream)
			for i := 0; i < limit; i++ {
				select {
				case <-done:
					return
				case outputStream <- <-intStream:
				}
			}
		}()
		return outputStream
	}

	primeFinder := func(done <-chan any, intStream <-chan int) <-chan int {
		primeStream := make(chan int)
		go func() {
			defer close(primeStream)
			for {
			Select:
				select {
				case <-done:
					return
				case i := <-intStream:
					copy := i
					for copy > 2 {
						copy -= 1
						if i%copy == 0 {
							break Select
						}
					}
					primeStream <- i
				}
			}
		}()
		return primeStream
	}

	fanInInt := func(done <-chan any, channels ...<-chan int) <-chan int {
		var wg sync.WaitGroup
		multiplexedStream := make(chan int)

		multiplex := func(c <-chan int) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	rand := func() any {
		return rand.Intn(50000000)
	}

	done := make(chan any)
	defer close(done)

	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	finders := make([]<-chan int, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range takeInt(done, fanInInt(done, finders...), 100) {
		fmt.Println(prime)
	}

}
