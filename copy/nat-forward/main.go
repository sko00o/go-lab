package main

import (
	"time"

	"github.com/sko00o/go-lab/copy/nat-forward/client"
	"github.com/sko00o/go-lab/copy/nat-forward/server"
)

func main() {
	go server.Run()
	time.Sleep(time.Millisecond)
	client.Run()
}
