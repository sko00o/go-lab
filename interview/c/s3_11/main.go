package main

import "time"

func main() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	ct := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}
		ct++
		time.Sleep(1 * time.Second)
	}

	println("done, ct=", ct)
}
