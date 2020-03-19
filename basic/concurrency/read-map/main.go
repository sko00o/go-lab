package main

import (
	"fmt"
	"sync"
)

func main() {
	// make an map
	m := make(map[int]int)

	// init data
	for i := 0; i < 100; i++ {
		m[i] = 100 - i
	}

	f := func(i int, wg *sync.WaitGroup) {
		defer wg.Done()
		v, ok := m[i]
		if !ok {
			fmt.Printf("invalid key; %d", i)
			return
		}

		fmt.Println(v)
	}

	var wg sync.WaitGroup
	gct := 10000
	wg.Add(gct)
	for g := 0; g < gct; g++ {
		go f(g%100, &wg)
	}

	wg.Wait()
}
