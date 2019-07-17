package mqv1

import (
	"sync"

	"github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const (
	ClientID   = "MQTEST1"
	CloseDelay = 250
)

type Backend struct {
	Client  mqtt.Client
	Topic   string
	OutChan chan mqtt.Message
	qos     byte
	mutex   sync.RWMutex
	logger  *log.Entry
}

type MqttSettings struct {
	Address      string `ini:"mqtt_tcp_hostname"`
	Username     string `ini:"mqtt_tcp_username"`
	Password     string `ini:"mqtt_tcp_password"`
	Qos          int    `ini:"mqtt_tcp_qos"`
	Poolsize     int    `ini:"mqtt_client_poolsize"`
	PublishSleep int64  `ini:"mqtt_publish_sleep"`
	Topic        string
}

func NewBackend(borker *MqttSettings) (*Backend, error) {
	b := &Backend{
		OutChan: make(chan mqtt.Message, 5000),
		qos:     byte(borker.Qos),
		logger:  log.WithField("name", "mq_old"),
		Topic:   borker.Topic,
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(borker.Address)
	opts.SetClientID(ClientID)
	opts.SetUsername(borker.Username)
	opts.SetPassword(borker.Password)
	opts.AutoReconnect = true
	opts.CleanSession = true
	opts.SetOnConnectHandler(b.onConnected)
	opts.SetConnectionLostHandler(b.onLost)

	b.Client = mqtt.NewClient(opts)
	if token := b.Client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return b, nil
}

func (b *Backend) Close() {
	b.Client.Disconnect(CloseDelay)
	b.logger.Info("closed")
}

func (b *Backend) onLost(client mqtt.Client, err error) {
	b.logger.Errorf("lost connection: %v", err)
}

func (b *Backend) onConnected(client mqtt.Client) {
	b.logger.Info("connected")
	if b.Client != client {
		b.logger.Infof("client: %v, b.Client: %v", client, b.Client)
	}
	if err := b.Subscribe(b.Topic); err != nil {
		b.logger.Errorf("subscrib error: %v", err)
	}
	// go func(t string) {
	// retry:
	// 	if token := b.Client.Subscribe(t, b.qos, b.handler); token.WaitTimeout(time.Duration(2)*time.Second) && token.Error() != nil {
	// 		log.Info(token.Error())
	// 		time.Sleep(time.Second)
	// 		goto retry
	// 	}
	// }(b.Topic)
}

func (b *Backend) Subscribe(topic string) error {
	defer b.mutex.Unlock()
	b.mutex.Lock()
	if token := b.Client.Subscribe(topic, b.qos, b.handler); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	b.logger.Infof("subscribing to topic %v", topic)
	return nil
}

func (b *Backend) UnSubscribe(topic string) error {
	defer b.mutex.Unlock()
	b.mutex.Lock()
	// log.Infof("backend/Backend: unsubscribing from topic %v", topic)
	if token := b.Client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (b *Backend) Publish(index int, topic string, bytes []byte) error {
	if t := b.Client.Publish(topic, b.qos, false, bytes); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}

func (b *Backend) handler(client mqtt.Client, msg mqtt.Message) {
	b.logger.Infoln(string(msg.Payload()))
}
