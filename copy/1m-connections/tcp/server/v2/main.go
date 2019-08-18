package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	ln, err := net.Listen("tcp", ":8964")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	epoller, err := MkEpoll()
	if err != nil {
		panic(err)
	}
	go start(epoller)

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
		if err := Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}
}

func start(epoller *epoll) {
	var buf = make([]byte, 8)
	for {
		connections, err := Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}

		for _, conn := range connections {
			if conn == nil {
				break
			}
			if _, err := conn.Read(buf); err != nil {
				if err := Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
				conn.Close()
			}
		}
	}
}
