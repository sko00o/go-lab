package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/sko00o/go-lab/echo/api"
	"github.com/sko00o/go-lab/echo/cert"
)

const (
	address     = "localhost:8686"
	defaultName = "world"
)

func main() {
	keyPair, certPool := cert.GetCert()
	_ = keyPair

	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(certPool, address)
	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoServiceClient(conn)

	// Contact the server and print out its response
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r1, err := c.Hello(context.Background(), &pb.EchoMessage{Body: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf(r1.Body)

	r2, err := c.Echo(context.Background(), &pb.EchoMessage{Body: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf(r2.Body)
}
