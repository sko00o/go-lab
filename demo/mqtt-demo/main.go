package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/sko00o/go-lab/demo/mqtt-demo/mqv1"
	"github.com/sko00o/go-lab/demo/mqtt-demo/mqv2"
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

func mq1() *mqv1.Backend {
	mq1, err := mqv1.NewBackend(&mqv1.MqttSettings{
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

func mq2() *mqv2.Backend {
	mq2 := mqv2.NewBackend("mqtt://logs:logs123@192.168.20.47:1883", topic)
	return mq2
}

func main() {
	app := cli.NewApp()
	app.Action = run

	app.Run(os.Args)
}
