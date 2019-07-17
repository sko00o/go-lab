package main

import "sync/atomic"

var value int32

func SetValue(delta int32) {
	for {
		v := value // 获取 value 不是并发安全的 改成 atomic.LoadInt32(&value)
		if atomic.CompareAndSwapInt32(&value, v, (v + delta)) {
			break
		}
	}
}

func main() {}
