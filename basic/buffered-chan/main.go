package main

import (
	"fmt"
	"sync"
	"time"
)

func gen(l int, ch chan int) {
	defer close(ch)

	var wg sync.WaitGroup
	wg.Add(l)
	for i := 0; i < l; i++ {
		time.Sleep(time.Millisecond)
		go func(idx int) {
			ch <- idx
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func cmp(ch chan int) {
	start := time.Now()
	go gen(10000, ch)

	// waiting for all goroutine blocked
	time.Sleep(11 * time.Second)

	cIdx := 0
	wrongNum := 0
	for idx := range ch {
		fmt.Println(idx, cIdx == idx)
		if cIdx != idx {
			wrongNum++
		}

		cIdx++
	}

	fmt.Printf("spent: %s, wrong: %d\n", time.Since(start), wrongNum)
}

func main() {
	ch := make(chan int)
	cmp(ch)

	ch1 := make(chan int, 10000)
	cmp(ch1)
}

/*
conclusion:
channel FIFO
*/
