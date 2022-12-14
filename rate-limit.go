package main

import (
	"context"
	"log"
	"os"
	"sync"

	"golang.org/x/time/rate"
)

type APIConnection struct {
	rateLimiter *rate.Limiter
}

func Open() *APIConnection {
	return &APIConnection{
		rateLimiter: rate.NewLimiter(rate.Limit(1), 1),
	}
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func (a *APIConnection) ResoleAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func RateLimit() {
	defer log.Printf("Done.")

	// per := func(eventCount int, duration time.Duration) rate.Limit {
	// 	return rate.Every(duration / time.Duration(eventCount))
	// }

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot read file: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResoleAddress(context.Background())
			if err != nil {
				log.Printf("cannot resolve address: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}

	wg.Wait()
}
