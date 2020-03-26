package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	sig := make(chan os.Signal)
	exit := make(chan struct{})
	signal.Notify(sig, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go run(ctx, exit)

	log.Printf("will quit by %v\n", <-sig)
	cancel()
	select {
	case <-exit:
	case s := <-sig:
		log.Printf("force quit by %v\n", s)
	}
}

func run(ctx context.Context, exit chan<- struct{}) {
	defer func() {
		exit <- struct{}{}
	}()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Topic:          "test_1",
		Brokers:        []string{"192.168.20.102:9092"},
		GroupID:        "test_group_1",
		MinBytes:       1,
		MaxBytes:       10e6,
		MaxWait:        100 * time.Millisecond,
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
	})

	go func() {
		<-ctx.Done()
		_ = reader.Close()
	}()

	log.Println("reading...")
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			if err == io.EOF {
				log.Printf("reader[%s] closed", reader.Config().Topic)
				return
			}

			log.Println("Err:", err)
			break
		}

		log.Println("recv", msg.Value)
	}

	log.Println("read over")
}
