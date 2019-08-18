package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v12 "github.com/sko00o/go-lab/copy/todo-list/pkg/api/v1"
	logger2 "github.com/sko00o/go-lab/copy/todo-list/pkg/logger"
	middleware2 "github.com/sko00o/go-lab/copy/todo-list/pkg/protocol/rest/middleware"
)

func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v12.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		logger2.Log.Fatal("failed to start HTTP gateway", zap.String("reason", err.Error()))
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: middleware2.AddRequestID(middleware2.AddLogger(logger2.Log, mux)),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	logger2.Log.Info("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}
