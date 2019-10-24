package server

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "[server] ", log.LstdFlags)

func Run() {
	go makeControl("127.0.0.1:8009")
	go makeAccept("127.0.0.1:8007") // out coming
	go makeForward("127.0.0.1:8008")
	go releaseConnMatch()
	tcpForward()
}
