package v1

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

	"github.com/sko00o/go-lab/copy/echo/api"
	cert2 "github.com/sko00o/go-lab/copy/echo/cert"
	service2 "github.com/sko00o/go-lab/copy/echo/pkg/service"
)

func makeGRPCServer(certPool *x509.CertPool, bind string) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(certPool, bind)),
	}

	// setup grpc server
	s := grpc.NewServer(opts...)
	api.RegisterEchoServiceServer(s, service2.NewServer())

	// register reflection service on gRPC server.
	reflection.Register(s)
	return s
}

func getRestMux(certPool *x509.CertPool, bind string, opts ...runtime.ServeMuxOption) (*runtime.ServeMux, error) {
	// we run out REST endpoint on the same port as the gRPC address.
	upstreamGRPCServerAddress := bind

	// get context, this allows control of the connection
	ctx := context.Background()

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: upstreamGRPCServerAddress,
		RootCAs:    certPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard,
		&runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	err := api.RegisterEchoServiceHandlerFromEndpoint(ctx, gwmux, upstreamGRPCServerAddress, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return nil, err
	}

	return gwmux, nil
}

// Note: ServeHTTP use http2, so TLS required
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func RunServer(bind string) {
	keyPair, certPool := cert2.GetCert()

	grpcServer := makeGRPCServer(certPool, bind)

	restMux, err := getRestMux(certPool, bind)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/foobar/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is a test endpoint"))
	})
	mux.Handle("/", restMux)

	mergeHandler := grpcHandlerFunc(grpcServer, mux)
	srv := &http.Server{
		Addr:    bind,
		Handler: mergeHandler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*keyPair},
			NextProtos:   []string{"h2"},
		},
	}

	conn, err := net.Listen("tcp", bind)
	if err != nil {
		panic(err)
	}

	fmt.Printf("starting GRPC and REST on: %v\n", bind)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
