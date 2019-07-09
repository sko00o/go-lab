package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/sko00o/go-lab/echo/api"
	"github.com/sko00o/go-lab/echo/pkg/service"
)

// join the two constants for convenience
var serveAddress = fmt.Sprintf("%v:%d", "localhost", 8686)

func simpleHTTPHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is a test endpoint"))
}

func makeGRPCServer() *grpc.Server {

	//setup grpc server
	s := grpc.NewServer()
	pb.RegisterEchoServiceServer(s, service.NewServer())
	// Register reflection service on gRPC server.
	reflection.Register(s)
	return s
}

func serveGRPC(l net.Listener) {

	s := makeGRPCServer()

	if err := s.Serve(l); err != nil {
		log.Fatalf("While serving gRpc request: %v", err)
	}
}

func serveHTTP(l net.Listener) {
	if err := http.Serve(l, nil); err != nil {
		log.Fatalf("While serving http request: %v", err)
	}
}

func main() {

	// Create a listener at the desired port.
	l, err := net.Listen("tcp", serveAddress)
	defer l.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux object.
	tcpm := cmux.New(l)

	// Declare the match for different services required.
	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	grpcL := tcpm.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := tcpm.Match(cmux.HTTP1Fast())

	// Link the endpoint to the handler function.
	// http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/foobar/", simpleHTTPHello)

	// Initialize the servers by passing in the custom listeners (sub-listeners).
	go serveGRPC(grpcL)
	go serveHTTP(httpL)

	log.Println("grpc server started.")
	log.Println("http server started.")
	log.Println("Server listening on ", serveAddress)

	// Start cmux serving.
	if err := tcpm.Serve(); !strings.Contains(err.Error(),
		"use of closed network connection") {
		log.Fatal(err)
	}
}