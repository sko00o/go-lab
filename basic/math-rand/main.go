package main

import (
	"fmt"
	"math/rand"
)

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
