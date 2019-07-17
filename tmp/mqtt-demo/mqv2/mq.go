package mqv2

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const (
	ClientID   = "MQTEST2"
	CloseDelay = 250
)

type Backend struct {
	Client mqtt.Client
}

func NewBackend(rawURL string, subTopic string) *Backend {
	uri, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	b := Backend{
		Client: connect(ClientID, uri, func(c mqtt.Client) {
			log.WithField("name", "mq_new").Info("connected")
			if t := c.Subscribe(subTopic, 2, testHandler); t.Wait() && t.Error() != nil {
				log.WithField("name", "mq_new").Errorf("subscrib error: %v", t.Error())
			}
			log.WithField("name", "mq_new").Infof("subscribing to topic %v", subTopic)
		}),
	}
	return &b
}

func (b *Backend) Close() {
	b.Client.Disconnect(CloseDelay)
	log.WithField("name", "mq_new").Info("closed")
}

func connect(clientID string, uri *url.URL, onConnHandler func(c mqtt.Client)) mqtt.Client {
	opts := createClientOptions(clientID, uri)
	opts.SetOnConnectHandler(onConnHandler)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return client
}

func createClientOptions(clientID string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	if password, isSet := uri.User.Password(); isSet {
		opts.SetPassword(password)
	}
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.WithField("name", "mq_new").Errorf("lost connection: %v", err)
	})

	opts.SetClientID(clientID)
	return opts
}

func testHandler(c mqtt.Client, msg mqtt.Message) {
	log.WithField("name", "mq_new").Infoln(string(msg.Payload()))
}

func (b *Backend) Send(topic string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if token := b.Client.Publish(topic, 2, false, bytes); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
