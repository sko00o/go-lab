package main

import (
	"fmt"
	"sync"
)

// map 没有初始化

type UserAges struct {
	ages map[string]int
	sync.Mutex
}

func (ua *UserAges) Add(name string, age int) {
	ua.Lock()
	defer ua.Unlock()

	// fix
	if ua.ages == nil {
		ua.ages = make(map[string]int)
	}

	ua.ages[name] = age
}

func (ua *UserAges) Get(name string) int {
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}

func main() {
	ua := new(UserAges)
	fmt.Println(ua.Get("a"))
	ua.Add("a", 1)
	fmt.Println(ua.Get("a"))
}
