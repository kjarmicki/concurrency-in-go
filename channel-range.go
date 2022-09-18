package main

import (
	"fmt"
	"time"
)

/*
 * Range operator can be used to iterate over channel values until the channel is closed
 */

func ChannelRange() {
	values := make(chan int)
	go func() {
		defer close(values)
		for i := 0; i < 5; i++ {
			values <- i
			time.Sleep(time.Second)
		}
	}()

	for val := range values {
		fmt.Println(val)
	}
}
