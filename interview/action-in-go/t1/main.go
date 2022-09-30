package main

import (
	"fmt"
	"time"
)

// 合并输出自然数

var count = 15

func ping(c chan<- int) {
	for i := 1; i < count; i++ {
		c <- 2*i - 1
	}
}
func pong(c chan<- int) {
	for i := 1; i < count; i++ {
		c <- 2 * i
	}
}

func print(ch <-chan int) {
	for {
		msg := <-ch
		fmt.Println(msg)
		time.Sleep(time.Millisecond * 50)
	}
}

func main() {
	ch := make(chan int)
	go ping(ch)
	go pong(ch)
	go print(ch)

	var input string
	fmt.Scanln(&input)
}
