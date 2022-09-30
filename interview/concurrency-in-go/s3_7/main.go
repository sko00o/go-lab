package main

import "sync"

// once for deadlock

func main() {
	var o1, o2 sync.Once

	var f2 func()
	f1 := func() { o2.Do(f2) }
	f2 = func() { o1.Do(f1) }

	f2() // can not exit
}
