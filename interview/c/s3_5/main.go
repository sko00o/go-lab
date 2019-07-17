package main

import (
	"fmt"
	"sync"
)

// use cond broadcast

func main() {
	type button struct {
		Clicked *sync.Cond
	}
	btn := button{sync.NewCond(&sync.Mutex{})}

	sub := func(c *sync.Cond, fn func()) {
		var run sync.WaitGroup
		run.Add(1)
		go func() {
			run.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		run.Wait()
	}

	var clickReg sync.WaitGroup
	clickReg.Add(3)
	sub(btn.Clicked, func() {
		fmt.Println("func 1")
		clickReg.Done()
	})
	sub(btn.Clicked, func() {
		fmt.Println("func 2")
		clickReg.Done()
	})
	sub(btn.Clicked, func() {
		fmt.Println("func 3")
		clickReg.Done()
	})

	btn.Clicked.Broadcast()
	clickReg.Wait()
}
