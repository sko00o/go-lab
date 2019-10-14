package main

import (
	"encoding/json"
	"fmt"
)

type J struct {
	data []int
}

// JSON 标准库对 nil slice 和 empty slice 处理不同。

func main() {
	out := func(label string, slice interface{}) {
		bs, err := json.Marshal(slice)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s:\n\t%s\n", label, string(bs))
	}

	var slice []int
	out("nil slice", slice)
	out("nil slice in struct", J{slice})

	slice = make([]int, 0)
	out("empty slice", slice)
	out("empty slice in struct", J{slice})
}
