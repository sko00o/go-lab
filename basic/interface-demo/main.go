package main

import "fmt"

func main() {
	t1 := newBase()
	fmt.Println(Talk())

	t2 := newHuman()
	fmt.Println(Talk(), Walk())

	t3 := newHuman21()
	fmt.Println(Talk(), Walk())

	t4 := newHuman22()
	fmt.Println(Talk(), Walk())
}
