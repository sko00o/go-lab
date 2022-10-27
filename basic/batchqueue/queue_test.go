package batchqueue

import (
	"sync"
	"testing"
	"time"
)

func TestQueue_Close(t *testing.T) {
	for _, tt := range []struct {
		name  string
		queue Queue
	}{
		{
			name:  "channel queue",
			queue: NewChannelQueue(5),
		},
		{
			name:  "slice queue",
			queue: NewSliceQueue(5),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					tt.queue.Get()
				}()
			}

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					tt.queue.Put('x')
				}()
			}

			time.Sleep(10 * time.Millisecond)
			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					tt.queue.Close()
				}()
			}

			wg.Wait()
		})
	}
}

func TestQueue(t *testing.T) {
	for _, tt := range []struct {
		name  string
		queue Queue
	}{
		{
			name:  "channel queue",
			queue: NewChannelQueue(5),
		},
		{
			name:  "slice queue",
			queue: NewSliceQueue(5),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			commonQueueTest(t, tt.queue)
		})
	}
}

func commonQueueTest(t *testing.T, q Queue) {
	var getWG sync.WaitGroup
	getWG.Add(1)
	go func() {
		defer getWG.Done()
		for {
			item := q.Get()
			if item == nil {
				t.Logf("nothing can get from queue")
				return
			}
			t.Logf("got: %v", item)
		}
	}()

	time.Sleep(1 * time.Second)

	itemCnt := 10
	t.Logf("concurrency put %d items", itemCnt)

	var putWG sync.WaitGroup
	for i := 0; i < itemCnt; i++ {
		item := i

		putWG.Add(1)
		go func() {
			defer putWG.Done()

			time.Sleep(time.Duration(item%2) * 100 * time.Millisecond)
			q.Put(item)
		}()
	}

	putWG.Wait()
	q.Close()
	getWG.Wait()

	if q.Get() != nil {
		t.Errorf("call Get() on closed queue, expect `nil`")
	}

	if q.Put(0) != false {
		t.Errorf("call Put() on closed queue, expect `false`")
	}

	// we can call Close more than once
	q.Close()
}
