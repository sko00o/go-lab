package main

import (
	"fmt"
	"sync"
	"time"
)

func deCnt(val int, w string, wg *sync.WaitGroup) {
	for i := 0; i < val; i++ {
		fmt.Println(w, i)
	}
	fmt.Println("over")
	wg.Done()
}

func main() {
	t1 := time.Now()
	val := 50000
	wg := sync.WaitGroup{}
	wg.Add(1)
	go deCnt(val, "go1: ", &wg)
	wg.Add(1)
	go deCnt(val, "go2: ", &wg)
	wg.Wait()
	fmt.Printf("time cost: %v", time.Now().Sub(t1))
}
