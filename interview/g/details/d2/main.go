package main

import "fmt"

func False() bool {
	return false
}

func main() {
	switch False(); { // 更好理解 d1 中的例子
	case false:
		fmt.Println("false")
	case true:
		fmt.Println("true")
	}

	switch False() {
	case false:
		fmt.Println("false")
	case true:
		fmt.Println("true")
	}
}
