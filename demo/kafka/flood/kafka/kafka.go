package kafka

import (
	"context"
	"io"
	"time"

	_ "github.com/segmentio/kafka-go/gzip"
	_ "github.com/segmentio/kafka-go/lz4"
	_ "github.com/segmentio/kafka-go/snappy"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Brokers []string
	GroupID string
	Topic   string
}

type Client struct {
	consumer *kafka.Reader
	producer *kafka.Writer
	ctx      context.Context
	cancel   context.CancelFunc
}

var dataRec = make(chan []byte)

var log = logrus.WithField("app", "kafka")

func Setup(c Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	// A consumer
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.Brokers,
		GroupID:     c.GroupID,
		Topic:       c.Topic,
		MinBytes:    1,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
	})

	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   c.Brokers,
		Topic:     c.Topic,
		BatchSize: 1,
		Balancer:  &kafka.Hash{},
	})

	go func(in chan []byte) {
		defer close(dataRec)
		for consumer != nil {
			msg, err := consumer.ReadMessage(ctx)
			if err != nil {
				if err == io.EOF || err == context.Canceled {
					break
				}

				log.WithError(err).Error("receive error")
				time.Sleep(2 * time.Second)
				continue
			}

			// check time
			dataRec <- msg.Value
		}
	}(dataRec)

	return &Client{
		consumer: consumer,
		producer: producer,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (c *Client) Send(data []byte) error {
	return c.producer.WriteMessages(c.ctx, kafka.Message{
		Key:   []byte("flood_data"),
		Value: data,
	})
}

func (c *Client) Recv() <-chan []byte {
	return dataRec
}

func (c *Client) Close() error {
	_ = c.producer.Close()
	_ = c.consumer.Close()
	c.cancel()

	return nil
}
