package main

import (
	"fmt"
)

type student struct {
	Name string
	Age  int
}

func makeStudent() {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}

	for _, stu := range stus {
		// cp := stu
		m[stu.Name] = &stu // range 是引用传递?
	}

	for _, stu := range m {
		fmt.Println(stu)
	}
}

func main() {
	makeStudent()
}
