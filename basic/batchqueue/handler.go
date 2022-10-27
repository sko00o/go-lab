package batchqueue

import (
	"sync"
	"time"
)

type Handler struct {
	// mutex protect the currBatch -u
	mutex     sync.Mutex
	currBatch *Batch

	// queue for input items
	queue Queue

	//
	wg sync.WaitGroup

	batchTimeout time.Duration
	batchSize    int
	batchBytes   int64

	handle handleFunc
}

type handleFunc func(batch *Batch)

func NewHandler(batchTimeout time.Duration, batchSize int, batchBytes int64, handle handleFunc, useChannelQueue bool) *Handler {
	h := &Handler{
		batchTimeout: batchTimeout,
		batchSize:    batchSize,
		batchBytes:   batchBytes,
		handle:       handle,
	}
	if useChannelQueue {
		h.queue = NewChannelQueue(10)
	} else {
		h.queue = NewSliceQueue(10)
	}

	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		h.handleLoop()
	}()

	return h
}

func (h *Handler) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.currBatch != nil {
		batch := h.currBatch
		h.queue.Put(batch)
		h.currBatch = nil
		batch.Fire()
	}
	h.queue.Close()
	h.wg.Wait()
}

func (h *Handler) Write(items ...Item) {
	h.wg.Add(1)
	defer h.wg.Done()

	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := range items {
	redo:
		batch := h.currBatch
		if batch == nil {
			batch = h.newBatch()
			h.currBatch = batch
		}
		if !batch.Add(items[i]) {
			batch.Fire()
			h.queue.Put(batch)
			h.currBatch = nil
			goto redo
		}

		if batch.Full() {
			batch.Fire()
			h.queue.Put(batch)
			h.currBatch = nil
		}
	}
}

func (h *Handler) handleLoop() {
	for {
		qi := h.queue.Get()
		if qi == nil {
			return
		}

		if batch, ok := qi.(*Batch); ok {
			if h.handle != nil {
				h.handle(batch)
			}
		}
	}
}

func (h *Handler) newBatch() *Batch {
	batch := newBatch(h.batchTimeout, h.batchSize, h.batchBytes)
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		h.awaitBatch(batch)
	}()
	return batch
}

func (h *Handler) awaitBatch(batch *Batch) {
	select {
	case <-batch.timer.C:
		h.mutex.Lock()
		if h.currBatch == batch {
			h.queue.Put(batch)
			h.currBatch = nil
		}
		h.mutex.Unlock()
	case <-batch.ready:
		batch.timer.Stop()
	}
}

type Batch struct {
	Items     []Item
	currSize  int
	currBytes int64
	maxSize   int
	maxBytes  int64
	ready     chan struct{}
	timer     *time.Timer
}

func newBatch(timeout time.Duration, maxSize int, maxBytes int64) *Batch {
	return &Batch{
		ready:    make(chan struct{}),
		timer:    time.NewTimer(timeout),
		maxSize:  maxSize,
		maxBytes: maxBytes,
	}
}

func (b *Batch) Add(item Item) bool {
	iBytes := item.Size()

	if b.currSize > 0 && b.currBytes+iBytes > b.maxBytes {
		return false
	}

	if cap(b.Items) == 0 {
		b.Items = make([]Item, 0, b.maxSize)
	}

	b.Items = append(b.Items, item)
	b.currSize++
	b.currBytes += iBytes
	return true
}

func (b *Batch) Full() bool {
	return b.currSize >= b.maxSize || b.currBytes >= b.maxBytes
}

func (b *Batch) Fire() {
	close(b.ready)
}
