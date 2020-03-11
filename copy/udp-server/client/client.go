package client

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

const maxBufferSize = 65507 // max udp data size

type Client struct {
	Ctx    context.Context
	Logger *zerolog.Logger
	Addr   string
	Reader io.Reader

	WillReceive bool
	Timeout     time.Duration
}

// Run implement a simple udp client
func (c *Client) Run() error {
	ctx, address, reader, logger := c.Ctx, c.Addr, c.Reader, c.Logger

	rAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			c.Logger.Err(err).Msg("conn close failed")
		}
	}()

	doneChan := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup
		wg.Add(2)

		// send
		go func() {
			defer wg.Done()

			_, err := io.Copy(conn, reader)
			if err != nil {
				doneChan <- err
				return
			}
		}()

		// receive
		go func() {
			defer wg.Done()

			if c.WillReceive {
				doneChan <- c.receive(ctx, conn, c.Timeout)
				return
			}
		}()

		wg.Wait()

		doneChan <- nil
	}()

	select {
	case <-ctx.Done():
		logger.Info().Msg("cancelled")
		return ctx.Err()

	case err := <-doneChan:
		return err
	}
}

func (c *Client) receive(ctx context.Context, conn *net.UDPConn, timeout time.Duration) error {
	for {
		deadline := time.Now().Add(timeout)
		err := conn.SetReadDeadline(deadline)
		if err != nil {
			return err
		}
		buffer := make([]byte, maxBufferSize)
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			// timeout error will not in log error
			if err, ok := err.(net.Error); ok && err.Timeout() {
				return nil
			}
			return err
		}
		c.Logger.Debug().Msgf("packet-received: bytes=%d from=%s", n, addr)

		c.Logger.Info().Msgf("receive: [%x]", buffer[:n])
	}
}
