package state

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func setUri(uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	//opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.AddBroker("tcp://" + uri.Host)
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	return opts
}

func (s *State) Connect(clientId string, uri *url.URL) mqtt.Client {
	opts := setUri(uri)
	logger := log.New(os.Stdout, "[Mqtt] ", log.LstdFlags)

	opts.SetClientID(clientId)

	// interval 2s
	opts.SetKeepAlive(2 * time.Second)
	opts.SetResumeSubs(true)

	// Lost callback
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logger.Println("Lost Connect")
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logger.Println("New Connect")

		// Disable mqtt v3 Sync status
		//err := s.Sync(client)
		//if err != nil {
		//	logger.Panicln("Sub error", err)
		//}
	})

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if err := token.Error(); err != nil {
		logger.Fatal(err)
	}

	s.Mqtt = client
	return client
}

func (s *State) Sync(client mqtt.Client) error {
	token := client.Subscribe("nodes/+/msg/+", 1, func(client mqtt.Client, msg mqtt.Message) {
		id := strings.Split(msg.Topic(), "/")[1]
		str := strings.Split(msg.Topic(), "/")[3]

		s.NodePut(id, str, msg.Payload())
	})

	if err := token.Error(); err != nil {
		return err
	}

	token = client.Subscribe("nodes/+/status", 1, func(client mqtt.Client, msg mqtt.Message) {
		id := strings.Split(msg.Topic(), "/")[1]
		s.SetNodeStatus(id, msg.Payload())
	})

	if err := token.Error(); err != nil {
		return err
	}

	return nil
}
