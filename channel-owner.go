package main

import (
	"fmt"
)

/*
 * When returning a channel from a function, a typical split of responsibility is that a function will write
 * values to that channel and close it. Therefore a read-only channel is returned because consumers shouldn't
 * be able to write values or close the channel themselves.
 */

func ChannelOwner() {
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i < 5; i++ {
				resultStream <- i
				fmt.Printf("sent %d\n", i)
			}
			fmt.Println("done sending")
		}()
		return resultStream
	}

	readStream := chanOwner()
	for val := range readStream {
		fmt.Printf("received %d\n", val)
	}
	fmt.Println("done receiving")
}
