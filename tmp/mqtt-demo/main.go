package main

import (
	"os"
	"os/signal"
	"syscall"

	v1 "tmp/mqtt-demo/mqv1"
	v2 "tmp/mqtt-demo/mqv2"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	topic = "testmq"
)

func run(c *cli.Context) error {

	mq1 := mq1()
	defer mq1.Close()

	mq := mq2()
	defer mq.Close()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Infof("signal received signal %v", <-sigChan)
	log.Warn("shutting down server")
	return nil
}

func mq1() *v1.Backend {
	mq1, err := v1.NewBackend(&v1.MqttSettings{
		Address:  "tcp://192.168.20.47:1883",
		Username: "niserver",
		Password: "niserver123@$",
		Qos:      2,
		Topic:    topic,
	})
	if err != nil {
		log.Fatal(err)
	}

	// mq1.Subscribe(topic)
	return mq1
}

func mq2() *v2.Backend {
	mq2 := v2.NewBackend("mqtt://logs:logs123@192.168.20.47:1883", topic)
	return mq2
}

func main() {
	app := cli.NewApp()
	app.Action = run

	app.Run(os.Args)
}
