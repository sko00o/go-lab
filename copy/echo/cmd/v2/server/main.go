package main

import (
	"fmt"

	v22 "github.com/sko00o/go-lab/copy/echo/pkg/protocol/v2"
)

// join the two constants for convenience
var serveAddress = fmt.Sprintf("%v:%d", "localhost", 8686)

func main() {
	v22.RunServer(serveAddress)
}
