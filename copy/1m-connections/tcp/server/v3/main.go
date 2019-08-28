package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/libp2p/go-reuseport"
	"github.com/rcrowley/go-metrics"

	"github.com/sko00o/go-lab/copy/1m-connections/tcp/epoll"
	"github.com/sko00o/go-lab/copy/1m-connections/tcp/limit"
	"github.com/sko00o/go-lab/copy/1m-connections/tcp/pool"
)

var (
	c = flag.Int("c", 10, "concurrency")
)

func main() {
	flag.Parse()
	limit.SetLimit()

	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	epoller, err := epoll.MkEpoll()
	if err != nil {
		panic(err)
	}
	workerPool := pool.NewPool(*c, 1000000, func(conn net.Conn) {
		if err := epoller.Remove(conn); err != nil {
			log.Printf("failed to remove %v", err)
		}
		conn.Close()
	})
	workerPool.Start()

	go start(epoller, workerPool)

	startEpoll(epoller)

	// wait stop signal
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop

	workerPool.Close()
}

func startEpoll(epoller *epoll.Epoll) {
	ln, err := reuseport.Listen("tcp", ":8964")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
				log.Printf("accept temp err: %v", nErr)
				continue
			}

			log.Printf("accept err: %v", err)
			return
		}

		// use epoll
		if err := epoller.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}
}

func start(epoller *epoll.Epoll, workerPool *pool.Pool) {
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}

		for _, conn := range connections {
			if conn == nil {
				break
			}

			workerPool.AddTask(conn)
		}
	}
}
