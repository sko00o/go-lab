package main

import (
	"sync"
)

type void struct{}

type threadSafeSet struct {
	sync.RWMutex
	s map[interface{}]void // this is a set
}

// 下面的代码有什么问题?
// 这里能找到源码：https://github.com/deckarep/golang-set/blob/master/threadsafe.go#L164
// 作者的说法：https://github.com/deckarep/golang-set/issues/50#issuecomment-312178265
// 问题出在，如果没有完全消费这个迭代函数返回的 chan， 会产生 goroutine 泄露
func (set *threadSafeSet) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		// 读锁，同时只有一个协程可以读取 chan 中的元素
		set.RLock()

		for elem := range set.s {
			ch <- elem
		}

		close(ch)
		set.RUnlock()
	}()
	return ch
}

func (set *threadSafeSet) Store(elem int) {
	set.RLock()
	set.s[elem] = void{}
	set.RUnlock()
}

func main() {
	var wg sync.WaitGroup

	// 多个协程调用迭代函数
	goroutines := 2
	wg.Add(goroutines)
	set := write(10)
	for i := 0; i < goroutines; i++ {
		go func(id int, set *threadSafeSet) {
			read(id, set)
			wg.Done()
		}(i, set)
	}
	wg.Wait()
}

func write(n int) *threadSafeSet {
	set := threadSafeSet{
		s: make(map[interface{}]void),
	}

	// 有一个协程不断向 set 中写
	go func() {
		for i := 0; i < n; i++ {
			set.Store(i)
		}
	}()

	return &set
}

func read(id int, set *threadSafeSet) {
	ch := set.Iter()

outer:
	for {
		select {
		case v, ok := <-ch:
			if ok {
				println(id, "read", v.(int))
			} else {
				println("closed")
				break outer
			}
		}
	}

	println("done")
}
