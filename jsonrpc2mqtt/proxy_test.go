package jsonrpc2mqtt

import (
	"testing"
)

func Test_mqttproxy(t *testing.T) {
	client := &MockClient{}

	mqttProxy, _ := OpenMqttProxy(client)
	mqttProxy.Notify("233", []byte("2222222222233333333333333"))

	go mqttProxy.SyncRpc("233", []byte(`{"jsonrpc":"2.0","method":"test","params":{"a":"233","b":"456"},"id":"test.0"}`))

	message := &MockMessage{
		DataPayload: []byte(`{"jsonrpc":"2.0","id":"test.0","result":"test2"}`),
	}
	client.SubscribeCallback(client, message)
	mqttProxy.SyncRpc("233", []byte(`{"jsonrpc":"2.0","id":"test.0","result":"test2"}`))
}
