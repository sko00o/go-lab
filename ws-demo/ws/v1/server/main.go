package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func main() {
	ln, err := net.Listen("tcp", ":8964")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v")
		}
	}()

	var connections []net.Conn
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

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

		go handleConn(conn)
		connections = append(connections, conn)
		if len(connections)%100 == 0 {
			log.Printf("total number of connetions: %v", len(connections))
		}
	}
}

func handleConn(conn net.Conn) {
	io.Copy(ioutil.Discard, conn)
}
