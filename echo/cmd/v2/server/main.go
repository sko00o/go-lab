package main

import (
	"fmt"

	v2 "github.com/sko00o/go-lab/echo/pkg/protocol/v2"
)

// join the two constants for convenience
var serveAddress = fmt.Sprintf("%v:%d", "localhost", 8686)

func main() {
	v2.RunServer(serveAddress)
}
