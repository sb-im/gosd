package jsonrpc2mqtt

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MockToken struct{}

func (t *MockToken) Wait() bool {
	return true
}

func (t *MockToken) WaitTimeout(time time.Duration) bool {
	return true
}

func (t *MockToken) Error() error {
	return errors.New("")
}

type MockClient struct {
	PublishData       []byte
	SubscribeCallback mqtt.MessageHandler
}

func (m *MockClient) IsConnected() bool {
	return true
}

func (m *MockClient) IsConnectionOpen() bool {
	return true
}

func (m *MockClient) Connect() mqtt.Token {
	return &MockToken{}
}

func (m *MockClient) Disconnect(quiesce uint) {}

func (m *MockClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	m.PublishData = payload.([]byte)
	return &MockToken{}
}

func (m *MockClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	m.SubscribeCallback = callback
	return &MockToken{}
}

func (m *MockClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return &MockToken{}
}
func (m *MockClient) Unsubscribe(topics ...string) mqtt.Token {
	return &MockToken{}
}
func (m *MockClient) AddRoute(topic string, callback mqtt.MessageHandler) {}

func (m *MockClient) OptionsReader() mqtt.ClientOptionsReader {
	return mqtt.ClientOptionsReader{}
}

type MockMessage struct {
	DataTopic   string
	DataPayload []byte
}

func (m *MockMessage) Duplicate() bool {
	return true
}

func (m *MockMessage) Qos() byte {
	return byte(0)
}

func (m *MockMessage) Retained() bool {
	return true
}

func (m *MockMessage) Topic() string {
	return m.DataTopic
}

func (m *MockMessage) MessageID() uint16 {
	return 233
}

func (m *MockMessage) Payload() []byte {
	return m.DataPayload
}

func (m *MockMessage) Ack() {}
