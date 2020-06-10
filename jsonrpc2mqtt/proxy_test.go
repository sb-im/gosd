package jsonrpc2mqtt

import (
	"errors"
	"testing"
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

type mqttClient struct {
	PublishData       []byte
	SubscribeCallback mqtt.MessageHandler
}

func (m *mqttClient) IsConnected() bool {
	return true
}

func (m *mqttClient) IsConnectionOpen() bool {
	return true
}

func (m *mqttClient) Connect() mqtt.Token {
	return &MockToken{}
}

func (m *mqttClient) Disconnect(quiesce uint) {}

func (m *mqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	m.PublishData = payload.([]byte)
	return &MockToken{}
}

func (m *mqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	m.SubscribeCallback = callback
	return &MockToken{}
}

func (m *mqttClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return &MockToken{}
}
func (m *mqttClient) Unsubscribe(topics ...string) mqtt.Token {
	return &MockToken{}
}
func (m *mqttClient) AddRoute(topic string, callback mqtt.MessageHandler) {}

func (m *mqttClient) OptionsReader() mqtt.ClientOptionsReader {
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

func Test_mqttproxy(t *testing.T) {
	client := &mqttClient{}

	mqttProxy, _ := OpenMqttProxy(client)
	mqttProxy.Notify("233", []byte("2222222222233333333333333"))

	go mqttProxy.SyncRpc("233", []byte(`{"jsonrpc":"2.0","method":"test","params":{"a":"233","b":"456"},"id":"test.0"}`))

	message := &MockMessage{
		DataPayload: []byte(`{"jsonrpc":"2.0","id":"test.0","result":"test2"}`),
	}
	client.SubscribeCallback(client, message)
	mqttProxy.SyncRpc("233", []byte(`{"jsonrpc":"2.0","id":"test.0","result":"test2"}`))
}
