package main

import "fmt"

var a [5]int
var p *[7]string

// 编译期用内置函数 len cap 赋值

// N and M are both typed constants.
const N = len(a)
const M = cap(p)

func main() {
	fmt.Println(N) // 5
	fmt.Println(M) // 7
}
