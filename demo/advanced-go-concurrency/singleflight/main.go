package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// 降低慢操作在高并发时的调用次数

func main() {
	start := time.Now()
	defer fmt.Println("cost:", time.Since(start))

	var group singleflight.Group
	var wg sync.WaitGroup
	for idx := 0; idx < 100000; idx++ {
		wg.Add(1)
		go groupCall(&wg, &group, idx)
	}
	wg.Wait()
}

func groupCall(wg *sync.WaitGroup, group *singleflight.Group, idx int) {
	defer wg.Done()
	_, _, _ = group.Do("call_slow_func", func() (i interface{}, e error) {
		slowFunc(idx)
		return
	})
}

func slowFunc(idx int) {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("call", idx)
}
