package multi_worker_decode

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
)

// 性能提升不大
func TestDecode(t *testing.T) {
	workers := 10
	times := 1000000

	// prepare the data
	c := Complex{
		A: Simple2{
			B: Simple1{
				C: Simple{A: "aaa"},
				B: 2,
			},
			C: &Simple1{
				C: Simple{A: "bbb"},
				B: 3,
			},
		},
		G: "abc",
		H: []int{1, 2, 3},
	}
	data, _ := json.Marshal(c)

	start := time.Now()
	out := make(chan Complex, 1)
	fn := prepare(workers, out)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var ct int
		for range out {
			ct++
			if ct == times {
				return
			}
		}
	}()
	for i := 0; i < times; i++ {
		fn(data)
	}
	wg.Wait()
	t.Log(time.Since(start))
}
