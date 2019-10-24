package server

import (
	"net"
	"time"
)

var cache *net.TCPConn

// 监听内网客户端的连接
func makeControl(addr string) {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", addr)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	logger.Println("listening: ", addr)

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		logger.Printf("control connected: %s <- %s\n", addr, tcpConn.RemoteAddr())

		if cache != nil {
			logger.Println("client already exist!")
			tcpConn.Close()
		} else {
			cache = tcpConn
		}
		go control(tcpConn)
	}
}

func control(conn *net.TCPConn) {
	go func() {
		for {
			// send "hi" every 2 seconds.
			_, err := conn.Write([]byte("hi\n"))
			if err != nil {
				cache = nil
			}
			time.Sleep(2 * time.Second)
		}
	}()
}
