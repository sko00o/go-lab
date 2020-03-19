package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// 加权限制并发

func main() {
	ctx := context.TODO()
	sem := semaphore.NewWeighted(10)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		cost := int64(i % 3)
		if err := sem.Acquire(ctx, cost); err != nil {
			break
		}

		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			defer sem.Release(cost)

			doSomething(idx, cost)
		}(i)
	}

	wg.Wait()
}

func doSomething(idx int, cost int64) {
	fmt.Printf("do %d cost %d\n", idx, cost)
	time.Sleep(time.Second)
}
