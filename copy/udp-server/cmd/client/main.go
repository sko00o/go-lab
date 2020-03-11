package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sko00o/go-lab/copy/udp-server/client"
)

var (
	address     = flag.String("d", "127.0.0.1:2333", "connect host:port")
	mode        = flag.String("m", "auto", "[input|auto]")
	willReceive = flag.Bool("r", false, "will receive")
	timeout     = flag.Duration("t", 5*time.Second, "receive timeout")

	debugMode = flag.Bool("debug", false, "enable debug mode")

	logger zerolog.Logger

	sig chan os.Signal
)

func init() {
	flag.Parse()

	sig = make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Str("component", "client").
		Logger()
	if *debugMode {
		logger = logger.Level(zerolog.DebugLevel)
	} else {
		logger = logger.Level(zerolog.InfoLevel)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := client.Client{
			Ctx:         ctx,
			Logger:      &logger,
			Addr:        *address,
			WillReceive: *willReceive,
			Timeout:     *timeout,
		}
		switch *mode {
		case "auto":
			c.Reader = autoGen()
		case "input":
			fmt.Println("type what you want, hit enter to go.")
			c.Reader = os.Stdin

		default:
			logger.Error().Msg("invalid mode")
		}

		if err := c.Run(); err != nil {
			logger.Err(err).Msg("client quit")
			return
		}
	}()

	select {
	case s := <-sig:
		logger.Info().Msgf("signal received: %v", s)
		cancel()
	}
}

type dd struct {
	idx int32
}

func (d *dd) Read(p []byte) (n int, err error) {
	binary.BigEndian.PutUint16(p, uint16(d.idx))
	d.idx = (d.idx + 1) % (1 << 16)
	return 2, nil
}

func autoGen() (r io.Reader) {
	return &dd{idx: 1}
}
