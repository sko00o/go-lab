package main

import "fmt"

func main() {
	t1 := newBase()
	fmt.Println(t1.Talk())

	t2 := newHuman()
	fmt.Println(t2.Talk(), t2.Walk())

	t3 := newHuman21()
	fmt.Println(t3.Talk(), t3.Walk())

	t31 := base(t3)
	fmt.Println(t31.Talk())

	t4 := newHuman22()
	fmt.Println(t4.Talk(), t4.Walk())
}
