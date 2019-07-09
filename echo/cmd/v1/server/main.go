package main

import (
	"fmt"

	v1 "github.com/sko00o/go-lab/echo/pkg/protocol/v1"
)

var serverAddress = fmt.Sprintf("%v:%d", "localhost", 8686)

func main() {
	v1.RunServer(serverAddress)
}
