package server

import (
	"net"
	"strconv"
	"sync"
	"time"
)

// 监听对外端口
func makeAccept(addr string) {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", addr)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			logger.Println("accept connect error: ", err)
			continue
		}

		logger.Printf("accept client connected %s <- %s\n", addr, tcpConn.RemoteAddr())
		addConnMatchAccept(tcpConn)
		sendMessage("new\n")
	}
}

type ConnMatch struct {
	accept        *net.TCPConn
	acceptAddTime time.Time
	tunnel        *net.TCPConn
}

var connListMap = make(map[string]*ConnMatch)

var lock = sync.Mutex{}

func addConnMatchAccept(accept *net.TCPConn) {
	lock.Lock()
	defer lock.Unlock()
	now := time.Now()
	connListMap[strconv.FormatInt(now.UnixNano(), 10)] = &ConnMatch{
		accept:        accept,
		acceptAddTime: now,
		tunnel:        nil,
	}
}

func sendMessage(msg string) {
	if cache == nil {
		logger.Println("send message error: no client")
	}

	_, err := cache.Write([]byte(msg))
	if err != nil {
		logger.Println("write message error: ", err)
		return
	}
	logger.Println("send message: ", msg)
}
