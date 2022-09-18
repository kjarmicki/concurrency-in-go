package main

import (
	"fmt"
	"time"
)

/*
 * Basically the same thing as in pattern-done-channel.go but applied to a loop.
 */

func PatternForSelectLoop_LoopAndWaitToBeStopped() {
	doneOwner := func() <-chan struct{} {
		done := make(chan struct{})
		go func() {
			<-time.After(time.Second)
			close(done)
		}()
		return done
	}

	done := doneOwner()
	var doWork func(i int)
	doWork = func(i int) {
		for {
			select {
			case <-done:
				return
			default:
			}

			fmt.Println(i)
			doWork(i + 1)
		}
	}

	doWork(0)
}
