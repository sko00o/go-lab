package service

import (
	"context"
	"log"

	"github.com/sko00o/go-lab/copy/echo/api"
)

type server struct{}

// implements hello function of EchoServiceServer
func (s *server) Hello(ctx context.Context, in *api.EchoMessage) (*api.EchoMessage, error) {
	log.Printf("Echo: get Hello")
	return &api.EchoMessage{Body: "Hello from server!"}, nil
}

// implements echo function of EchoServiceServer
func (s *server) Echo(ctx context.Context, in *api.EchoMessage) (*api.EchoMessage, error) {
	log.Printf("Echo: Received '%v'", in.Body)
	return &api.EchoMessage{Body: "ACK " + in.Body}, nil
}

func NewServer() api.EchoServiceServer {
	return &server{}
}
