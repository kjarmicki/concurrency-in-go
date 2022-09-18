package main

import (
	"log"
	"os"
	"time"
)

func HealingGoroutines() {
	type startGoroutineFn func(
		done <-chan any,
		pulseInterval time.Duration,
	) (heartbeat <-chan any)

	either := func(this <-chan any, that <-chan any) <-chan any {
		merged := make(chan any)
		go func() {
			defer close(merged)

			select {
			case <-this:
			case <-that:
			}
		}()
		return merged
	}

	newSteward := func(
		timeout time.Duration,
		startGoroutine startGoroutineFn,
	) startGoroutineFn {
		return func(done <-chan any, pulseInterval time.Duration) <-chan any {
			heartbeat := make(chan any)
			go func() {
				defer close(heartbeat)
				var wardDone chan any
				var wardHeartbeat <-chan any
				startWard := func() {
					wardDone = make(chan any)
					wardHeartbeat = startGoroutine(either(done, wardDone), timeout/2)
				}
				startWard()
				pulse := time.Tick(pulseInterval)

			monitorLoop:
				for {
					timeoutSignal := time.After(timeout)

					for {
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case <-wardHeartbeat:
							continue monitorLoop
						case <-timeoutSignal:
							log.Println("steward: ward unhealthy, restarting")
							close(wardDone)
							startWard()
							continue monitorLoop
						case <-done:
							return
						}
					}
				}
			}()
			return heartbeat
		}
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <-chan any, _ time.Duration) <-chan any {
		go func() {
			log.Println("ward: started")
			<-done
			log.Println("ward: ended")
		}()
		return nil
	}
	doWorkWithSteward := newSteward(4*time.Second, doWork)

	done := make(chan any)
	time.AfterFunc(9*time.Second, func() {
		log.Println("main: halting steward and ward")
		close(done)
	})

	for range doWorkWithSteward(done, 4*time.Second) {
	}
	log.Println("done")
}
