package main

import (
	"fmt"
	"time"
)

/*
 * A worker function accepts two channels: one with values to be processed
 * and the other to receive a signal of when to stop working.
 */

func PatternDoneChannel() {
	doWork := func(
		done <-chan any,
		strings <-chan string,
	) <-chan any {
		terminated := make(chan any)
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan any)

	fmt.Println("starting work")
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(time.Second)
		fmt.Println("Canceling doWork goroutine")
		close(done)
	}()

	<-terminated
	fmt.Println("done")
}
