package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb "github.com/sko00o/go-lab/echo/api"
	"github.com/sko00o/go-lab/echo/cert"
)

var serverAddress = fmt.Sprintf("%v:%d", "localhost", 8686)

type server struct{}

// implements Hello func of EchoServiceServer
func (s *server) Hello(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	log.Printf("Echo: get Hello")
	return &pb.EchoMessage{Body: "Hello from server"}, nil
}

// implements Echo func of EchoServiceServer
func (s *server) Echo(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	log.Printf("Echo: Received '%v'", in.Body)
	return &pb.EchoMessage{Body: "ACK " + in.Body}, nil
}

func simpleHTTPHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is a test endpoint"))
}

func makeGRPCServer(certPool *x509.CertPool) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(certPool, serverAddress)),
	}

	// setup grpc server
	s := grpc.NewServer(opts...)
	pb.RegisterEchoServiceServer(s, &server{})

	// register reflection service on gRPC server.
	reflection.Register(s)
	return s
}

func getRestMux(certPool *x509.CertPool, opts ...runtime.ServeMuxOption) (*runtime.ServeMux, error) {
	// we run out REST endpoint on the same port as the gRPC address.
	upstreamGRPCServerAddress := serverAddress

	// get context, this allows control of the connection
	ctx := context.Background()

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: upstreamGRPCServerAddress,
		RootCAs:    certPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard,
		&runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	err := pb.RegisterEchoServiceHandlerFromEndpoint(ctx, gwmux, upstreamGRPCServerAddress, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return nil, err
	}

	return gwmux, nil
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	keyPair, certPool := cert.GetCert()

	grpcServer := makeGRPCServer(certPool)
	restMux, err := getRestMux(certPool)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/foobar/", simpleHTTPHello)
	mux.Handle("/", restMux)

	mergeHandler := grpcHandlerFunc(grpcServer, mux)
	srv := &http.Server{
		Addr:    serverAddress,
		Handler: mergeHandler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*keyPair},
			NextProtos:   []string{"h2"},
		},
	}

	conn, err := net.Listen("tcp", serverAddress)
	if err != nil {
		panic(err)
	}

	fmt.Printf("starting GRPC and REST on: %v\n", serverAddress)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
