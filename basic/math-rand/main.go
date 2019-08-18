package main

import (
	"fmt"
	"math/rand"
)

// fake random
func main() {
	data := make([]byte, 4)
	rand.Read(data)
	fmt.Println(data)
	rand.Read(data)
	fmt.Println(data)
	rand.Read(data)
	fmt.Println(data)
	rand.Read(data)
	fmt.Println(data)
}

/*

[82 253 252 7]
[33 130 101 79]
[22 63 95 15]
[154 98 29 114]

*/
