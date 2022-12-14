package main

import (
	"fmt"
	"runtime"
	"sync"
)

// goroutine resource counting

const numGoroutines = 1e6

func main() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c }

	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}

/*
numGoutines = 1e4

run on linux: go1.12, go1.12.4 linux/amd64
0.151kb

run on windows: go1.11.5, go1.12.5 windows/amd64
8.718kb
*/

/*
numGoutines = 1e6

run on linux: go1.12, go1.12.4 linux/amd64
2.584kb

run on windows: go1.11.5, go1.12.5 windows/amd64
8.7958kb

run on macOS: go1.19.1 darwin/arm64
2.620kb
*/
