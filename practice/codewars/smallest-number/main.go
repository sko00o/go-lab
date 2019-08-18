package main

import (
	"fmt"
	"strconv"
)

// From codewars
// 给定一个数
// 移动某一位数
// 得到最小的数

func Smallest(n int64) []int64 {
	res := []int64{n, 0, 0}
	s := strconv.FormatInt(n, 10)
	for i := range s {
		s2 := s[:i] + s[i+1:]
		for j := range s {
			s3 := s2[:j] + string(s[i]) + s2[j:]
			n2, err := strconv.ParseInt(s3, 10, 64)
			if err == nil && n2 < res[0] {
				res = []int64{n2, int64(i), int64(j)}
			}
		}
	}
	return res
}

func main() {

	fmt.Println(Smallest(935855753))
}
