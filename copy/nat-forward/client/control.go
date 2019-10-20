package client

import (
	"bufio"
	"io"
	"net"
)

func connControl(control string, local, remote string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", control)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Println("control connect error: ", err)
		return
	}
	logger.Printf("control connected: %s -> %s\n", conn.LocalAddr(), control)

	reader := bufio.NewReader(conn)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				logger.Println("read control error:", err)
			}
			break
		}

		switch s {
		case "new\n":
			go combine(local, remote)
		case "hi\n":
			// ignore
		}
	}
}

func combine(localAddr, remoteAddr string) {
	local := connectLocal(localAddr)
	remote := connectRemote(remoteAddr)
	if local != nil && remote != nil {
		joinConn(local, remote)
	} else {
		if local != nil {
			if err := local.Close(); err != nil {
				logger.Println("close local error: ", err)
			}
			if err := remote.Close(); err != nil {
				logger.Println("close remote error: ", err)
			}
		}
	}
}

func joinConn(local, remote *net.TCPConn) {
	f := func(local *net.TCPConn, remote *net.TCPConn) {
		defer local.Close()
		defer remote.Close()
		_, err := io.Copy(local, remote)
		if err != nil {
			if err != io.EOF {
				logger.Println("io copy error:", err)
			}
			return
		}
		logger.Println("end")
	}
	go f(local, remote)
	go f(remote, local)
}

func connectLocal(addr string) *net.TCPConn {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Println("local connect error: ", err)
		return nil
	}

	logger.Printf("local connected: %s -> %s\n", conn.LocalAddr(), addr)
	return conn
}

func connectRemote(addr string) *net.TCPConn {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Println("remote connect error: ", err)
		return nil
	}

	logger.Printf("remote connected: %s -> %s\n", conn.LocalAddr(), addr)
	return conn
}
