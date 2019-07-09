package service

import (
	"context"
	"log"

	pb "github.com/sko00o/go-lab/echo/api"
)

type server struct{}

// implements hello function of EchoServiceServer
func (s *server) Hello(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	log.Printf("Echo: get Hello")
	return &pb.EchoMessage{Body: "Hello from server!"}, nil
}

// implements echo function of EchoServiceServer
func (s *server) Echo(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	log.Printf("Echo: Received '%v'", in.Body)
	return &pb.EchoMessage{Body: "ACK " + in.Body}, nil
}

func NewServer() pb.EchoServiceServer {
	return &server{}
}
