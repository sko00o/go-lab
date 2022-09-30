package srv

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func warmServiceCache() *sync.Pool {
	p := &sync.Pool{
		New: connectToService,
	}

	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func startNetworkDaemon(bind string) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		srv, err := net.Listen("tcp", bind)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer srv.Close()

		wg.Done()

		for {
			conn, err := srv.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			connectToService()
			fmt.Fprintln(conn, "")
			conn.Close()
		}
	}()

	return &wg
}

func startNetworkDaemon2(bind string) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceCache()

		srv, err := net.Listen("tcp", bind)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer srv.Close()

		wg.Done()

		for {
			conn, err := srv.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			srvConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(srvConn)
			conn.Close()
		}
	}()

	return &wg
}
