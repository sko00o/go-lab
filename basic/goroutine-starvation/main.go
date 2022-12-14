package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	done := make(chan bool, 1)
	var mu sync.Mutex

	start := time.Now()
	var g1, g2 int

	// goroutine 1 (greedy)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mu.Lock()
				g1++
				time.Sleep(10 * time.Microsecond)
				mu.Unlock()
			}
		}
	}()

	// goroutine 2
	for i := 0; i < 10; i++ {
		time.Sleep(10 * time.Microsecond)
		mu.Lock()
		g2++
		mu.Unlock()
	}
	done <- true

	fmt.Printf("Lock acquired per goroutine:\ng1: %d, g2: %d\ncost: %v\n", g1, g2, time.Since(start))
}
