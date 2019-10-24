package client

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "[client] ", log.LstdFlags)

func Run() {
	connControl(
		"127.0.0.1:8009",
		"127.0.0.1:8000", // behind nat
		"127.0.0.1:8008",
	)
}
