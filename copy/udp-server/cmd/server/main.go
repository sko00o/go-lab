package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sko00o/go-lab/copy/udp-server/server"
)

var (
	address  = flag.String("d", "0.0.0.0:2333", "listen host:port")
	sendBack = flag.Bool("b", false, "will send back")
	timeout  = flag.Duration("t", 5*time.Second, "send timeout")

	logger zerolog.Logger

	sig chan os.Signal
)

func init() {
	flag.Parse()

	sig = make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Str("component", "server").
		Logger()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := server.Run(server.Options{
			Ctx:      ctx,
			Logger:   &logger,
			Addr:     *address,
			SendBack: *sendBack,
			Timeout:  *timeout,
		}); err != nil {
			logger.Fatal().Err(err)
		}
	}()

	select {
	case s := <-sig:
		logger.Info().Msgf("signal received: %v", s)
		cancel()
	}
}
