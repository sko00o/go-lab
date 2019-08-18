package main

import (
	"fmt"

	v12 "github.com/sko00o/go-lab/copy/echo/pkg/protocol/v1"
)

var serverAddress = fmt.Sprintf("%v:%d", "localhost", 8686)

func main() {
	v12.RunServer(serverAddress)
}
