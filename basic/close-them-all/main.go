package main

import (
	"context"
	"sync"
)

func main() {
	count := 100
	ctx, cancel := context.WithCancel(context.Background())

	var start sync.WaitGroup
	start.Add(count)

	var end sync.WaitGroup
	end.Add(count + 1)

	for i := 0; i < count; i++ {
		go func(i int) {
			defer end.Done()
			start.Done()

			select {
			case <-ctx.Done():
				println("done", i)
				break
			}
		}(i)
	}

	go func() {
		defer end.Done()
		select {
		case <-ctx.Done():
			println("I am free")
		}
	}()

	//start.Wait()
	cancel()
	end.Wait()
}
