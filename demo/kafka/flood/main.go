package main

import (
	"context"
	"encoding/binary"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"

	kafka2 "github.com/sko00o/go-lab/demo/kafka/flood/kafka"
)

var log = logrus.WithField("app", "flood")

var (
	Brokers = []string{"192.168.20.102:9092"}
	GroupID = "flood"
	Topic   = "flood_test"
)

var (
	// create a cache storage
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

type MQ interface {
	Send(data []byte) error
	Recv() <-chan []byte
	Close() error
}

func main() {
	var client MQ

	// init kafka client
	client = kafka2.Setup(kafka2.Config{
		Brokers: Brokers,
		GroupID: GroupID,
		Topic:   Topic,
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go runConsumer(&wg, client)

	wg.Add(1)
	pCtx, pCancel := context.WithCancel(context.Background())
	go runProducer(pCtx, &wg, client)

	// wait for quit signal
	sig := make(chan os.Signal)
	exit := make(chan struct{})
	signal.Notify(sig, os.Interrupt)
	log.WithField("signal", <-sig).Info("signal received")
	go func() {
		pCancel()
		_ = client.Close()
		wg.Wait()
		exit <- struct{}{}
	}()
	select {
	case <-exit:
	case <-sig:
		log.WithField("signal", <-sig).Info("stop immediately")
	}
}

func runConsumer(wg *sync.WaitGroup, mq MQ) {
	defer wg.Done()

	for d := range mq.Recv() {
		data := make([]byte, len(d))
		copy(data, d)
		go func(data []byte) {
			iStart, ok := storage.Get(string(data))
			if !ok {
				log.Error("data not exist")
				return
			}
			start, ok := iStart.(time.Time)
			if !ok {
				log.Error("data invalid")
				return
			}
			log.Infof("data: %x, time cost: %v", data, time.Since(start))
		}(data)
	}
}

func runProducer(ctx context.Context, wg *sync.WaitGroup, mq MQ) {
	defer wg.Done()

	process := func(data []byte) error {
		if err := mq.Send(data); err != nil {
			return err
		}

		// save data
		storage.Set(string(data), time.Now(), cache.DefaultExpiration)
		return nil
	}

	// data generator
	t := time.NewTicker(1 * time.Millisecond)
	go func() {
		select {
		case <-ctx.Done():
			t.Stop()
		}
	}()
	i := 1
	for range t.C {
		i %= 1 << 16
		data := make([]byte, 2)
		binary.BigEndian.PutUint16(data[:], uint16(i))
		if err := process(data); err != nil {
			log.WithError(err).Error("send error")
			continue
		}
		i++
	}
}
