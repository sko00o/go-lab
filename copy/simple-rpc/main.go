package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

type User struct {
	Name string
	Age  int
}

var userDB = map[int]User{
	1: {"AAA", 12},
	9: {"BBB", 24},
	7: {"CCC", 36},
}

func QueryUser(id int) (User, error) {
	if u, ok := userDB[id]; ok {
		return u, nil
	}

	return User{}, fmt.Errorf("id %d not in user db", id)
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("error: %v", p)
		}
	}()

	gob.Register(User{})
	addr := "localhost:8964"
	srv := NewServer(addr)

	// start server
	srv.Register("QueryUser", QueryUser)
	go srv.Run()

	// wait for server to start
	time.Sleep(2 * time.Second)

	// start client
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	cli := NewClient(conn)

	var Query func(int) (User, error)
	cli.callRPC("QueryUser", &Query)

	u, err := Query(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	u2, err := Query(2)
	if err != nil {
		panic(err)
	}
	fmt.Println(u2)
}
