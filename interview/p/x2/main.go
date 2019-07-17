package main

import (
	"fmt"
)

// select buffered channel 有关的坑

var ct = 0

func httpHandler() {
	// errCh := make(chan error,1) // 有 buffer 的情况，会进入并发 select，进入异常流程
	// resultCh := make(chan int,1)
	errCh := make(chan error)
	resultCh := make(chan int)
	go func() {
		defer close(errCh)
		defer close(resultCh)
		errCh <- fmt.Errorf("shit") // 增加了 errCh 被 capture 的概率
	}()

	select {
	case <-errCh:
	case <-resultCh:
		println("this shall not happen dude!", ct) // 这条有可能被输出，可能因为 select 的 goroutine 正好轮询到这个
		ct++
	}
}

func main() {
	for i := 0; i < 1000000; i++ {
		httpHandler()
	}
}
