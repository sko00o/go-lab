package main

import (
	"fmt"
	"sync"
)

func main() {
	var s []int
	s = make([]int, 10)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s[i] = i
			// s = append(s, i) // WRONG!
		}(i)
	}

	wg.Wait()
	fmt.Printf("len: %d, cap: %d, %v", len(s), cap(s), s)
}
