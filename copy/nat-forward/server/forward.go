package server

import (
	"io"
	"net"
	"time"
)

func makeForward(addr string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()
	logger.Println("forward listening: ", addr)

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			logger.Println(err)
			continue
		}
		logger.Printf("forward connected: %s <- %s\n", addr, tcpConn.RemoteAddr())
		configureListenTunnel(tcpConn)
	}
}

var connListMapUpdate = make(chan int)

const (
	UPDATE = iota
)

func configureListenTunnel(tunnel *net.TCPConn) {
	lock.Lock()
	used := false

	for _, connMatch := range connListMap {
		if connMatch.tunnel == nil && connMatch.accept != nil {
			connMatch.tunnel = tunnel
			used = true
			break
		}
	}

	if !used {
		logger.Println("map size: ", len(connListMap))
		_ = tunnel.Close()
		logger.Println("tunnel closed")
	}
	lock.Unlock()
	connListMapUpdate <- UPDATE
}

func tcpForward() {
	for {
		select {
		case <-connListMapUpdate:
			lock.Lock()
			for key, connMatch := range connListMap {
				if connMatch.tunnel != nil && connMatch.accept != nil {
					logger.Println("new connection")
					go joinConn(connMatch.accept, connMatch.tunnel)
					delete(connListMap, key)
				}
			}
			lock.Unlock()
		}
	}
}

func joinConn(conn1, conn2 *net.TCPConn) {
	f := func(local, remote *net.TCPConn) {
		defer local.Close()
		defer remote.Close()
		_, err := io.Copy(local, remote)
		if err != nil {
			if err != io.EOF {
				logger.Println("io copy error:", err)
			}
			return
		}
		logger.Println("join conn end")
	}

	go f(conn2, conn1)
	go f(conn1, conn2)
}

func releaseConnMatch() {
	for {
		lock.Lock()
		for key, connMatch := range connListMap {
			if connMatch.tunnel == nil && connMatch.accept != nil {
				if time.Since(connMatch.acceptAddTime) > 5*time.Second {
					logger.Println("release connection")
					err := connMatch.accept.Close()
					if err != nil {
						logger.Println("release error: ", err)
					}
					delete(connListMap, key)
				}
			}
		}
		lock.Unlock()
		time.Sleep(5 * time.Second)
	}
}
