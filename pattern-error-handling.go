package main

import (
	"fmt"
	"net/http"
)

/*
 * "Returning" error from goroutines - create a result struct that holds both response and error.
 */

func PatternErrorHandling() {
	type Result struct {
		Response *http.Response
		Error    error
	}

	checkStatus := func(done <-chan any, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)

			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{
					Response: resp,
					Error:    err,
				}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan any)
	defer close(done)

	urls := []string{"https://google.com", "badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %s\n", result.Error.Error())
			continue
		}
		fmt.Printf("response status: %d\n", result.Response.StatusCode)
	}
}
