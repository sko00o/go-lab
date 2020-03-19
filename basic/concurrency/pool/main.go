package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	poolCount := 10
	jobChan := make(chan int, poolCount)
	for i := 0; i < poolCount; i++ {
		go func(idx int) {
			for j := range jobChan {
				fmt.Printf("job: %d, got: %d\n", idx, j)
				wg.Done()
			}
		}(i)
	}

	jobCount := 100
	for i := 0; i < jobCount; i++ {
		wg.Add(1)
		jobChan <- i
	}

	wg.Wait()
}
