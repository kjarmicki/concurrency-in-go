package main

import "fmt"

/*
 * Lexical confinement - using lexical scope to confine access to the channel.
 * Basically the same thing as in channel-owner.go but with an inline function
 */

func PatternConfinement() {
	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i < 10; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Println(result)
		}
	}

	consumer(chanOwner())
}
