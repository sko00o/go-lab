// +build !linux

package epoll

import (
	"net"
	"sync"
)

type Epoll struct {
	fd          int
	connections map[int]net.Conn
	lock        *sync.RWMutex
}

func MkEpoll() (*Epoll, error) {

	return nil, nil
}

func (e *Epoll) Add(conn net.Conn) error {
	return nil
}

func (e *Epoll) Remove(conn net.Conn) error {
	return nil
}

func (e *Epoll) Wait() ([]net.Conn, error) {

	return nil, nil
}
