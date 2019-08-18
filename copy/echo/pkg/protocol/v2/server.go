package v2

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sko00o/go-lab/copy/echo/api"
	service2 "github.com/sko00o/go-lab/copy/echo/pkg/service"
)

func serveGRPC(l net.Listener) {
	//setup grpc server
	s := grpc.NewServer()
	api.RegisterEchoServiceServer(s, service2.NewServer())
	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(l); err != nil {
		log.Fatalf("While serving gRpc request: %v", err)
	}
}

func serveHTTP(l net.Listener) {
	if err := http.Serve(l, nil); err != nil {
		log.Fatalf("While serving http request: %v", err)
	}
}

func RunServer(bind string) {
	// Create a listener at the desired port.
	l, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// Create a cmux object.
	tcpm := cmux.New(l)

	// Declare the match for different services required.
	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	grpcL := tcpm.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := tcpm.Match(cmux.HTTP1Fast())

	// Link the endpoint to the handler function.
	// http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/foobar", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is a test endpoint"))
	})

	// Initialize the servers by passing in the custom listeners (sub-listeners).
	go serveGRPC(grpcL)
	go serveHTTP(httpL)

	log.Println("grpc server started.")
	log.Println("http server started.")
	log.Println("Server listening on ", bind)

	// Start cmux serving.
	if err := tcpm.Serve(); !strings.Contains(err.Error(),
		"use of closed network connection") {
		log.Fatal(err)
	}
}
