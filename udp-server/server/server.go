package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	timeout   = flag.Duration("t", 5*time.Second, "send timeout")
	address   = flag.String("d", "0.0.0.0:2333", "listen host:port")
	sendBack  = flag.Bool("b", false, "will send back")
	debugMode = flag.Bool("debug", false, "enable debug mode")

	logger zerolog.Logger

	sig chan os.Signal
)

const maxBufferSize = 65507 // max udp data size

func init() {
	flag.Parse()

	sig = make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	if !*debugMode {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
			Level(zerolog.InfoLevel)
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("component", "server").Logger()
}

// Server implement a simple udp server
func Server(ctx context.Context, address string) error {

	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	defer pc.Close()

	doneChan := make(chan error, 1)

	buffer := make([]byte, maxBufferSize)

	logger.Info().Msg("server start")

	go func() {
		for {
			// receive
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}
			logger.Debug().Msgf("packet-received: bytes=%d from=%s", n, addr)

			// process
			logger.Info().Msgf("receive: [%x]", buffer[:n])

			// send back
			if *sendBack {
				deadline := time.Now().Add(*timeout)
				err = pc.SetWriteDeadline(deadline)
				if err != nil {
					doneChan <- err
					return
				}
				n, err = pc.WriteTo(buffer[:n], addr)
				if err != nil {
					doneChan <- err
					return
				}
				logger.Debug().Msgf("packet-written: bytes=%d to=%s", n, addr)
			}

		}
	}()

	select {
	case <-ctx.Done():
		logger.Info().Msg("cancelled")
		return ctx.Err()

	case err := <-doneChan:
		return err
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go Server(ctx, *address)

	select {
	case s := <-sig:
		log.Info().Msgf("signal received: %v", s)
		cancel()
	}
}
