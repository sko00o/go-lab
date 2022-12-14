package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestPool(t *testing.T) {
	p := &sync.Pool{
		New: func() interface{} {
			println("create a new instance.")
			return struct{}{}
		},
	}

	p.Get()
	ins := p.Get()
	p.Put(ins)
	p.Get()
}

func TestPool2(t *testing.T) {
	var numCalCreate int64
	calPool := &sync.Pool{
		New: func() interface{} {
			atomic.AddInt64(&numCalCreate, 1)
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// calPool.Put(calPool.New())
	// calPool.Put(calPool.New())
	// calPool.Put(calPool.New())
	// calPool.Put(calPool.New())

	const numWorker = 1 << 20
	var wg sync.WaitGroup
	wg.Add(numWorker)
	for i := numWorker; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calPool.Get().(*[]byte) // !
			defer calPool.Put(mem)
		}()
	}

	wg.Wait()
	println("number of calculator:", numCalCreate)
}
