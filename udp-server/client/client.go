package main

import (
	"context"
	"encoding/binary"
	"flag"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	mode        = flag.String("m", "auto", "[input|auto]")
	timeout     = flag.Duration("t", 5*time.Second, "receive timeout")
	address     = flag.String("d", "127.0.0.1:2333", "connect host:port")
	willReceive = flag.Bool("r", false, "will receive")
	debugMode   = flag.Bool("debug", false, "enable debug mode")

	logger zerolog.Logger

	sig chan os.Signal
)

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

	logger = log.With().Str("component", "client").Logger()
}

const maxBufferSize = 65507 // max udp data size

// Client implement a simple udp client
func Client(ctx context.Context, address string, reader io.Reader) error {

	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	doneChan := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup
		wg.Add(2)

		// send
		go func() {
			defer wg.Done()

			n, err := io.Copy(conn, reader)
			if err != nil {
				doneChan <- err
				return
			}
			logger.Debug().Msgf("packet-written: bytes=%d", n)
		}()

		// receive
		go func() {
			defer wg.Done()

			if *willReceive {
				doneChan <- reveive(ctx, conn)
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

	case s := <-sig:
		log.Info().Msgf("signal received: %v", s)
		return nil
	}
}

func loopClient(ctx context.Context, address string) error {

	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	doneChan := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup
		wg.Add(2)

		// send
		go func() {
			defer wg.Done()

			dt := make([]byte, 2)
			for i := 0; i < 1<<8; i++ {
				binary.BigEndian.PutUint16(dt, uint16(i))

				n, err := conn.Write(dt)
				if err != nil {
					logger.Error().Msgf("write error: %s", err)
					continue
				}
				logger.Debug().Msgf("packet-written: bytes=%d", n)

				logger.Info().Msgf("send: [%x]", dt)
			}
		}()

		// receive
		go func() {
			defer wg.Done()

			if *willReceive {
				doneChan <- reveive(ctx, conn)
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

	case s := <-sig:
		log.Info().Msgf("signal received: %v", s)
		return nil
	}
}

func reveive(ctx context.Context, conn *net.UDPConn) error {
	for {
		deadline := time.Now().Add(*timeout)
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
		logger.Debug().Msgf("packet-received: bytes=%d from=%s", n, addr)

		logger.Info().Msgf("receive: [%x]", buffer[:n])
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	switch *mode {
	case "auto":
		if err := loopClient(ctx, *address); err != nil {
			log.Error().Msg(err.Error())
		}

	case "input":
		if err := Client(ctx, *address, os.Stdin); err != nil {
			log.Error().Msg(err.Error())
		}

	default:
		logger.Error().Msg("invalid mode")
	}

	cancel()
}
