package batchqueue

import (
	"context"
	"sync"
)

type Queue interface {
	Put(item interface{}) bool
	Get() interface{}
	Close()
}

type ChannelQueue struct {
	ch chan interface{}

	ctx    context.Context
	cancel context.CancelFunc

	once sync.Once
}

func NewChannelQueue(fixedSize int) *ChannelQueue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &ChannelQueue{
		ch:     make(chan interface{}, fixedSize),
		ctx:    ctx,
		cancel: cancel,
	}
	return q
}

func (q *ChannelQueue) Put(item interface{}) bool {
	select {
	case <-q.ctx.Done():
		return false
	default:
		q.ch <- item
		return true
	}
}

func (q *ChannelQueue) Get() interface{} {
	item, ok := <-q.ch
	if !ok {
		return nil
	}
	return item
}

func (q *ChannelQueue) Close() {
	q.cancel()
	q.once.Do(func() {
		close(q.ch)
	})
}

// SliceQueue is a blocking queue implement on slice
type SliceQueue struct {
	items []interface{}
	cond  *sync.Cond

	closed bool
}

func NewSliceQueue(initSize int) *SliceQueue {
	q := &SliceQueue{
		items: make([]interface{}, 0, initSize),
		cond:  sync.NewCond(new(sync.Mutex)),
	}

	return q
}

func (q *SliceQueue) Put(item interface{}) bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	defer q.cond.Broadcast()

	if q.closed {
		return false
	}
	q.items = append(q.items, item)
	return true
}

func (q *SliceQueue) Get() interface{} {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.items) == 0 && !q.closed {
		q.cond.Wait()
	}

	if len(q.items) == 0 {
		return nil
	}

	item := q.items[0]
	q.items[0] = nil
	q.items = q.items[1:]

	return item
}

func (q *SliceQueue) Close() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	defer q.cond.Broadcast()

	q.closed = true
}
