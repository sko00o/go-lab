package main

import "fmt"

var a = &[]int{1, 2, 3}
var i int

func f() int {
	i = 1
	a = &[]int{7, 8, 9}
	return 0
}

func main() {
	// The evaluation order of "a", "i"
	// and "f()" is unspecified.
	(*a)[i] = f()
	fmt.Println(*a)
}
