package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func PoolNetwork() {
	connectToService := func() any {
		time.Sleep(time.Second)
		return struct{}{}
	}

	warmServiceConnCache := func() *sync.Pool {
		p := &sync.Pool{
			New: connectToService,
		}
		for i := 0; i < 10; i++ {
			p.Put(p.New())
		}
		return p
	}

	startNetworkDaemon := func() *sync.WaitGroup {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			pool := warmServiceConnCache()
			server, err := net.Listen("tcp", "localhost:8080")
			if err != nil {
				log.Fatalf("cannot listen: %v", err)
			}
			defer server.Close()
			wg.Done()

			for {
				conn, err := server.Accept()
				if err != nil {
					log.Printf("cannot accept connection: %v", err)
					continue
				}
				connect := pool.Get()
				fmt.Fprintln(conn, "")
				pool.Put(connect)
				conn.Close()
			}
		}()

		return &wg
	}

	wg := startNetworkDaemon()
	wg.Wait()
}
