package main

import "sync"

// use once

func main() {
	var ct int
	inc := func() {
		ct++
	}
	dec := func() {
		ct--
	}

	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(inc)
		}()
	}

	increments.Wait()
	println(ct)

	once.Do(dec) // Once only counts the number of times Do is called
	println(ct)
}
