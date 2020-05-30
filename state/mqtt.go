package state

import (
	"log"
	"net/url"
	"os"
	"strconv"
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
	})

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if err := token.Error(); err != nil {
		logger.Fatal(err)
	}

	client.Subscribe("nodes/+/msg/+", 1, func(client mqtt.Client, msg mqtt.Message) {
		m := string(msg.Payload())
		logger.Println("Recv:", strings.Split(msg.Topic(), "/")[3], m)

		id, _ := strconv.Atoi(strings.Split(msg.Topic(), "/")[2])
		str := strings.Split(msg.Topic(), "/")[3]

		if len(s.Node[id].Msg) == 0 {
			s.Node[id] = NodeState{
				Msg: map[string][]byte{},
			}
		}

		s.Node[id].Msg[str] = msg.Payload()
	})

	s.Mqtt = client

	return client
}
