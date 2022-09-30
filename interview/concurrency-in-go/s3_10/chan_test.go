package block

import (
	"sync"
	"testing"
)

func TestUnblockWithChan(t *testing.T) {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			t.Log(i, "has began")
		}(i)
	}

	t.Log("unblock all goroutine")
	close(begin)
	wg.Wait()
}

func TestUnblockWithCond(t *testing.T) {
	begin := sync.NewCond(&sync.Mutex{})
	var wg, enter sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)

		enter.Add(1)
		go func(i int) {
			defer wg.Done()
			begin.L.Lock()
			defer begin.L.Unlock()
			enter.Done()
			begin.Wait()
			println(i, "has began")
		}(i)
		// make sure all goroutine on Wait, otherwise Broadcast will work incorrect.
		enter.Wait()
	}

	println("unblock all goroutine")
	begin.Broadcast()
	wg.Wait()
}

func TestUnblockWithCondManyTimes(t *testing.T) {
	begin := sync.NewCond(&sync.Mutex{})
	var wg, enter sync.WaitGroup
	const callTimes = 3
	const routineNum = 5

	enter.Add(routineNum)
	for i := 0; i < routineNum; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			begin.L.Lock()
			defer begin.L.Unlock()
			for j := 0; j < callTimes; j++ {
				enter.Done() // TODO: need atomic!
				begin.Wait()
				t.Log(i, "has began for", j, "times")
			}
		}(i)

	}

	t.Log("unblock all goroutine")
	for i := 0; i < callTimes; i++ {
		enter.Wait()
		begin.Broadcast()
		enter.Add(routineNum)
	}
	wg.Wait()
}
