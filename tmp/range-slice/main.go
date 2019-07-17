package main

import (
	"bytes"
	"fmt"
)

func main() {
	tcs := []SS{
		{
			s: S{
				val: V{0, 1, 2},
			},
		},
		{
			s: S{
				val: V{2, 3, 4},
			},
		},
		{
			s: S{
				val: V{4, 5, 6},
			},
		},
	}

	for i, n := range tcs {
		fmt.Printf("%p, %p\n", &tcs[i], &n)
	}

	val := V{2, 3, 4}
	res := findRXInfo(tcs, val)
	fmt.Printf("%p : %v\n", res, res)

	val = V{0, 2, 4}
	res = findRXInfo(tcs, val)
	fmt.Printf("%p : %v\n", res, res)
}

type V [3]byte

type S struct {
	val V
}
type SS struct {
	s S
}

func findRXInfo(slice []SS, val V) *S {
	// for _, n := range slice // n 地址不变，迭代指向被遍历的切片
	for i := range slice {
		if bytes.Equal(slice[i].s.val[:], val[:]) {
			return &(slice[i].s)
		}
	}
	return nil
}
