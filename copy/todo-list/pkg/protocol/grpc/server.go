package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	v12 "github.com/sko00o/go-lab/copy/todo-list/pkg/api/v1"
	logger2 "github.com/sko00o/go-lab/copy/todo-list/pkg/logger"
	middleware2 "github.com/sko00o/go-lab/copy/todo-list/pkg/protocol/grpc/middleware"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, v1API v12.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// gRPC server startup options
	opts := []grpc.ServerOption{}
	// add middleware
	opts = middleware2.AddLogging(logger2.Log, opts)

	// register service
	server := grpc.NewServer(opts...)
	v12.RegisterToDoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger2.Log.Warn("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger2.Log.Info("starting gRPC server...")
	return server.Serve(listen)
}
