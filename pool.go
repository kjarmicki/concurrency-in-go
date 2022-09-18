package main

import (
	"fmt"
	"sync"
)

func UsePool() {
	myPool := &sync.Pool{
		New: func() any {
			fmt.Println("Creating new instance")
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}
