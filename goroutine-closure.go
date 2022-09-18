package main

import (
	"fmt"
	"sync"
)

/*
 * Go functions (and therefore also goroutines) are closures. This means that they have access to the enclosing scope,
 * but the variable value is read at the time of goroutine execution, not at the time of goroutine declaration.
 */

func GroroutineClosureWrong() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

func GroroutineClosureRight() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}
