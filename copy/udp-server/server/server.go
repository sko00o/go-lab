package server

import (
	"context"
	"net"
	"time"

	"github.com/rs/zerolog"
)

const maxBufferSize = 65507 // max udp data size

type Options struct {
	Ctx      context.Context
	Logger   *zerolog.Logger
	Addr     string
	SendBack bool
	Timeout  time.Duration
}

// Run implement a simple udp server
func Run(o Options) error {
	ctx, address, logger := o.Ctx, o.Addr, o.Logger

	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Err(err).Msg("conn close failed")
		}
	}()

	doneChan := make(chan error, 1)

	buffer := make([]byte, maxBufferSize)

	logger.Info().Msg("server start")

	go func() {
		for {
			// receive
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			// process
			logger.Info().
				Int("size", n).
				Str("from", addr.String()).
				Hex("data", buffer[:n]).
				Msg("RX")

			// send back
			if o.SendBack {
				go func(conn net.PacketConn, addr net.Addr, timeout time.Duration) {
					deadline := time.Now().Add(timeout)
					err = conn.SetWriteDeadline(deadline)
					if err != nil {
						doneChan <- err
						return
					}
					n, err = conn.WriteTo(buffer[:n], addr)
					if err != nil {
						doneChan <- err
						return
					}
					logger.Info().
						Int("size", n).
						Str("to", addr.String()).
						Msg("TX")
				}(conn, addr, o.Timeout)
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
