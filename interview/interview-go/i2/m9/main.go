package main

import "fmt"

type Student struct {
	name string
}

func main() {
	m := map[string]Student{"people": {"zhoujielun"}}
	// m["people"].name = "wuyanzu" // 不可以直接操作map中的成员参数
	m["people"] = Student{"wuyanzu"}

	fmt.Print(m["people"])
}
