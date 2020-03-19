package main

import (
	"fmt"
	"sync"
)

var task = []func(){
	rangeQuit,
}

func main() {
	for _, t := range task {
		t()
	}
}

// quit goroutine by range
func rangeQuit() {
	var wg sync.WaitGroup

	wg.Add(1)
	in := make(chan int)
	go func(in chan int) {
		defer func() {
			wg.Done()
			println("range quit out.")
		}()
		for i := range in {
			fmt.Println(i)
		}
	}(in)

	// quit the goroutine
	close(in)
	wg.Wait()
}
