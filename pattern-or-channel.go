package main

import (
	"fmt"
	"time"
)

/*
 * `or` function returns as soon as one of the channels from variadic argument sends a value
 */

func PatternOrChannel() {
	var or func(channels ...<-chan any) <-chan any
	or = func(channels ...<-chan any) <-chan any {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan any)
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan any {
		ch := make(chan any)
		go func() {
			defer close(ch)
			time.Sleep(after)
		}()
		return ch
	}

	start := time.Now()
	<-or(
		sig(time.Second*6),
		sig(time.Second*5),
		sig(time.Second*4),
		sig(time.Second*3),
		sig(time.Second*2),
		sig(time.Second*1),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
