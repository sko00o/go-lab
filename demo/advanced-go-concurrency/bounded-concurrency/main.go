package main

import (
	"fmt"
	"sync"
	"time"
)

// 有界并发： 使用信号量限制并发数

func main() {
	semaphore := make(chan struct{}, 10)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		acquire(semaphore)

		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			defer release(semaphore)

			doSomething(idx)
		}(i)
	}

	wg.Wait()
}

func acquire(sem chan<- struct{}) {
	sem <- struct{}{}
}

func release(sem <-chan struct{}) {
	<-sem
}

func doSomething(idx int) {
	fmt.Println("do", idx)
	time.Sleep(time.Second)
}
