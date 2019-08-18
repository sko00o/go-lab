package main

import (
	"fmt"
)

func main() {
	array()
	slice()
}

func array() {
	nums := [...]int{1, 2, 3, 4, 5, 6}

	maxIndex := len(nums) - 1
	for i, e := range nums { // range 只在 for 语句开始执行时求值一次； 求值结果会被复制
		if i == maxIndex {
			nums[0] += e
		} else {
			nums[i+1] += e
		}
	}

	fmt.Println(nums) // [7 3 5 7 9 11]
}

func slice() {
	nums := []int{1, 2, 3, 4, 5, 6}

	maxIndex := len(nums) - 1
	for i, e := range nums {
		if i == maxIndex {
			nums[0] += e
		} else {
			nums[i+1] += e
		}
	}

	fmt.Println(nums) // [22 3 6 10 15 21]
}
