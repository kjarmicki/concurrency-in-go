package main

import (
	"fmt"
	"sync"
)

/*
 * When a channel is closed, the reads from it won't block. We can use this to simulate events.
 * Here all goroutines will block (wait for an event) and then resume when the "begin" channel is closed.
 */

func ChannelCloseEvent() {
	begin := make(chan any)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%d has begun\n", i)
		}(i)
	}

	fmt.Println("beginning goroutines")
	close(begin)
	wg.Wait()
}
