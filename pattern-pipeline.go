package main

import "fmt"

/*
 * Pipeline pattern is useful for building stream processing.
 * Work is split with an initial channel and then passed through independent steps in the pipeline.
 * Notice the pattern of building a pipeline step:
 *
 *	pipeline := func(done <-chan any, inputStream <-chan) <-chan int {
 *		outputStream := make(chan any)
 *		go func() {
 *			defer close(outputStream)
 *			for _, i := range inputStream {
 *				select {
 *				case <-done:
 *					return
 *				case outputStream <- (COMPUTATION HERE):
 *				}
 *			}
 *		}()
 *		return outputStream
 *	}
 */

func PatternPipeline() {
	generator := func(done <-chan any, integers ...int) <-chan int {
		intStream := make(chan int, len(integers))
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(done <-chan any, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan any, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan any)
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}