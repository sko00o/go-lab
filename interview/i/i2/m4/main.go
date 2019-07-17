package main

import "fmt"

type P struct {
	N string
}

func (p *P) String() string {
	// runtime: goroutine stack exceeds 1000000000-byte limit
	// fatal error: stack overflow
	// 递归调用了 String()
	// return fmt.Sprintf("print: %v", p)
	return fmt.Sprintf("print: %v", *p)
}

func main() {
	p := &P{N: "a"}
	print(p.String())
}
