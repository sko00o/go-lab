package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	addr        = flag.String("ip", "127.0.0.1:8964", "Server IP")
	connections = flag.Int("conn", 1, "number of tcp connections")
)

func main() {
	flag.Parse()

	log.Printf("Connecting to %s", *addr)
	var conns []net.Conn
	for i := 0; i < *connections; i++ {
		c, err := net.DialTimeout("tcp", *addr, 10*time.Second)
		if err != nil {
			fmt.Println("failed to connect", i, err)
			i--
			continue
		}
		conns = append(conns, c)
		time.Sleep(time.Millisecond)
	}

	defer func() {
		for _, c := range conns {
			c.Close()
		}
	}()

	log.Printf("initial connect %d", len(conns))
	tts := time.Second
	if *connections > 100 {
		tts = time.Millisecond * 5
	}

	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			conn.Write([]byte("hello world\n\n"))
		}
	}
}
