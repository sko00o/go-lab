package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	jobCh := genJob(10)
	retCh := make(chan string)

	var wg sync.WaitGroup
	workPool(10, jobCh, retCh, &wg)

	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for ret := range retCh {
			time.Sleep(1)
			fmt.Println(ret)
		}
	}()

	wg.Wait()
	close(retCh)
	wg1.Wait()
}

func workPool(n int, jobCh <-chan int, retCh chan<- string, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(i, jobCh, retCh, wg)
	}
}

func worker(id int, jobCh <-chan int, retCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobCh {
		time.Sleep(time.Second)
		ret := fmt.Sprintf("worker %d processed job %d", id, job)
		retCh <- ret
	}
}

func genJob(n int) <-chan int {
	jobCh := make(chan int)
	go func() {
		defer close(jobCh)
		for i := 0; i < n; i++ {
			jobCh <- i
		}
	}()

	return jobCh
}
