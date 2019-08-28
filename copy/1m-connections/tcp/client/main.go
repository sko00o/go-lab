package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/rcrowley/go-metrics"

	"github.com/sko00o/go-lab/copy/1m-connections/tcp/epoll"
	"github.com/sko00o/go-lab/copy/1m-connections/tcp/limit"
)

var (
	addr        = flag.String("ip", "127.0.0.1:8964", "Server IP")
	connections = flag.Int("conn", 1, "number of tcp connections")
	startMetric = flag.String("sm", time.Now().Format(time.RFC3339), "start time point of all clients")
	c           = flag.Int("c", 100, "currency count")
)

var (
	opsRate = metrics.NewRegisteredTimer("ops", nil)
)

func main() {
	flag.Parse()
	limit.SetLimit()

	go func() {
		startPoint, err := time.Parse(time.RFC3339, *startMetric)
		if err != nil {
			panic(err)
		}
		time.Sleep(startPoint.Sub(time.Now()))

		metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	}()

	log.Printf("Connecting to %s", *addr)

	// multi epoller
	for i := 0; i < *c; i++ {
		go mkClient(*addr, *connections/(*c))
	}

	// wait stop signal
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func mkClient(addr string, connections int) {
	epoller, err := epoll.MkEpoll()
	if err != nil {
		panic(err)
	}
	go start(epoller)

	var conns []net.Conn
	for i := 0; i < connections; i++ {
		c, err := net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			fmt.Println("failed to connect", i, err)
			i--
			continue
		}

		// 全部连接丢入 epoll
		if err := epoller.Add(c); err != nil {
			log.Printf("failed to add connection %v", err)
			c.Close()
		}

		conns = append(conns, c)
	}

	log.Printf("initial connect %d", len(conns))
	tts := time.Second
	if connections > 100 {
		tts = time.Millisecond * 5
	}

	for i := range conns {
		time.Sleep(tts)

		writeTimestamp(conns[i], epoller)
	}
}

func start(epoller *epoll.Epoll) {
	var nano int64
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

			// 拿到服务端返回的时间戳
			if err := binary.Read(conn, binary.BigEndian, &nano); err != nil {
				log.Printf("failed to read %v", err)
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}

				conn.Close()
				continue
			} else {
				opsRate.Update(time.Duration(time.Now().UnixNano() - nano))
			}

			// 再次写回
			writeTimestamp(conn, epoller)
		}
	}
}

func writeTimestamp(conn net.Conn, epoller *epoll.Epoll) {
	if err := binary.Write(conn, binary.BigEndian, time.Now().UnixNano()); err != nil {
		log.Printf("failed to write timestamp %v", err)
		if err := epoller.Remove(conn); err != nil {
			log.Printf("failed to remove %v", err)
		}
	}
}
