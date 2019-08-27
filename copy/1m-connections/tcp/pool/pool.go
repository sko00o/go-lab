package pool

import (
	"io"
	"net"
	"sync"

	"github.com/rcrowley/go-metrics"
)

var (
	opsRate = metrics.NewRegisteredMeter("ops", nil)
)

type CallbackFunc func(conn net.Conn)

type Pool struct {
	workers   int
	maxTasks  int
	taskQueue chan net.Conn

	mu     sync.Mutex
	closed bool
	done   chan struct{}

	// callback when handle got error
	callback CallbackFunc
}

func NewPool(w int, t int, f CallbackFunc) *Pool {
	return &Pool{
		workers:   w,
		maxTasks:  t,
		taskQueue: make(chan net.Conn, t),
		done:      make(chan struct{}),

		callback: f,
	}
}

func (p *Pool) Close() {
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

func (p *Pool) AddTask(conn net.Conn) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}

	p.mu.Unlock()
	p.taskQueue <- conn
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		go p.startWorker()
	}
}

func (p *Pool) startWorker() {
	for {
		select {
		case <-p.done:
			return
		case conn := <-p.taskQueue:
			if conn != nil {
				handleConn(conn, p.callback)
			}
		}
	}
}

func handleConn(conn net.Conn, f CallbackFunc) {
	_, err := io.CopyN(conn, conn, 8)
	if err != nil {
		f(conn)
	}

	opsRate.Mark(1)
}
