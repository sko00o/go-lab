package main

import "fmt"

// switch 的默认布尔值
func main() {
	switch { // <=> switch true {
	case false:
		fmt.Println("false")
	case true:
		fmt.Println("true")
	}
}
