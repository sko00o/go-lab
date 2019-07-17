package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan int)

	start := time.Now()
	go func(l int) {
		var wg sync.WaitGroup
		wg.Add(l)
		defer close(ch)
		for i := 0; i < l; i++ {
			time.Sleep(time.Millisecond)
			go func(idx int) {
				ch <- idx
				wg.Done()
			}(i)
		}
		wg.Wait()
	}(10000)

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

/*
结论：
无论是否有 buffer
channel 的消费都是随机的
*/
