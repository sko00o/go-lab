package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	SUCCESS = "[O]"
	FAILED  = "[X]"
	WARN    = "[!]"
	LINKED  = "[+]"
	SEND    = "[>]"
	RECV    = "[<]"
)

var (
	servers string
	clients string
)

func init() {
	flag.StringVar(&servers, "s", "", "-s server1,server2")
	flag.StringVar(&clients, "c", "127.0.0.1:2333,127.0.0.1:8964", "-c client1,client2")
	flag.Parse()
}

func main() {
	ss := strings.Split(servers, ",")
	cs := strings.Split(clients, ",")
	var err error
	var ls1, ls2 net.Listener
	var c1, c2 net.Conn
	switch {
	case len(ss) >= 2:
		if len(ss) > 2 {
			log.Println(WARN, "will ignore", ss[2:])
		}

		ls1 = listen(ss[0])
		ls2 = listen(ss[1])

		for {
			conn1, err := ls1.Accept()
			if err != nil {
				errLog("accept connect error", err)
				continue
			}

			conn2, err := ls2.Accept()
			if err != nil {
				errLog("accept connect error", err)
				continue
			}

			forward(conn1, conn2)
		}

	case len(ss) == 1 && ss[0] != "":
		if len(cs) < 1 {
			log.Fatal(FAILED, "require at least 1 client")
		}
		ls1 = listen(ss[0])

		if len(cs) > 1 {
			log.Println(WARN, "will ignore", cs[1:])
		}

		for {
			c1, err = dial(cs[0])
			if err != nil {
				errLog("dial error will retry in 2 second", err)
				time.Sleep(2 * time.Second)
				continue
			}

			conn1, err := ls1.Accept()
			if err != nil {
				errLog("accept connect error", err)
				continue
			}

			forward(conn1, c1)
		}
	default:
		if len(cs) < 2 {
			log.Fatal(FAILED, "require at least 2 client or server")
		}

		for {
			c1, err = dial(cs[0])
			if err != nil {
				errLog("dial error will retry in 2 second", err)
				time.Sleep(2 * time.Second)
				continue
			}

			c2, err = dial(cs[1])
			if err != nil {
				errLog("dial error will retry in 2 second", err)
				time.Sleep(2 * time.Second)
				continue
			}

			forward(c1, c2)
		}
	}
}

func listen(bind string) net.Listener {
	tcpAddr := resolveAddr(bind)

	ls, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf(FAILED+" listen to [%s] error: %v", tcpAddr, err)
	}

	log.Println(LINKED, "listening on", tcpAddr)

	return ls
}

func dial(addr string) (net.Conn, error) {
	tcpAddr := resolveAddr(addr)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("dial for [%s] error: %w", tcpAddr, err)
	}

	log.Println(LINKED, "connected to", tcpAddr)

	return conn, nil
}

func resolveAddr(addr string) *net.TCPAddr {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(FAILED, "resolve address error:", err)
	}
	return tcpAddr
}

// forward will block
func forward(conn1, conn2 net.Conn) {
	log.Printf(LINKED+" transmit between [%s],[%s] <-> [%s],[%s]",
		conn1.LocalAddr(), conn1.RemoteAddr(), conn2.LocalAddr(), conn2.RemoteAddr())
	var wg sync.WaitGroup
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	wg.Wait()
}

func connCopy(dst, src net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		_ = src.Close()
		_ = dst.Close()
	}()

	_, err := io.Copy(dst, src)
	if err != nil {
		/*if err != io.EOF {
			errLog("io copy error", err)
		}*/
		return
	}
}

func errLog(msg string, err error) {
	log.Printf(FAILED+" %s: %v", msg, err)
}
