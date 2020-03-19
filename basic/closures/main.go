package main

import (
	"fmt"
)

var (
	disableFuncsA bool
)

var funcsA = groupFuncsB(
	f1,
	groupFuncsA(
		f2,
	),
	f3,
)

var tasks = []func(){
	f0,
	funcsA,
	f4,
}

func main() {
	run()
	disableFuncsA = true
	run()
}

func run() {
	for _, t := range tasks {
		t()
	}
}

func groupFuncsB(funcs ...func()) func() {
	return func() {
		for _, f := range funcs {
			f()
		}
	}
}

func groupFuncsA(funcs ...func()) func() {
	return func() {
		if !disableFuncsA {
			for _, f := range funcs {
				f()
			}
		}
	}
}

func f0() {
	fmt.Println("f0")
}

func f1() {
	fmt.Println("f1")
}

func f2() {
	fmt.Println("f2")
}

func f3() {
	fmt.Println("f3")
}

func f4() {
	fmt.Println("f4")
}
