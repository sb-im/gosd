package jsonrpc2mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsonrpc2 "github.com/sb-im/jsonrpc2"
)

const (
	rpc_topic   = "nodes/%s/jsonrpc"
	rpc_timeont = 1 * time.Hour
)

type MqttProxy struct {
	Client  mqtt.Client
	mutex   sync.Mutex // protects pending
	pending map[jsonrpc2.ID]chan []byte
}

func OpenMqttProxy(client mqtt.Client) (*MqttProxy, error) {
	mqttProxy := &MqttProxy{
		Client:  client,
		pending: make(map[jsonrpc2.ID]chan []byte),
	}

	token := mqttProxy.Client.Subscribe(fmt.Sprintf(rpc_topic, "+"), 1, func(client mqtt.Client, msg mqtt.Message) {
		jsonrpc_res := jsonrpc2.Jsonrpc{}
		err := json.Unmarshal(msg.Payload(), &jsonrpc_res)
		if err != nil || jsonrpc_res.ID == nil {
			return
		}

		if pending := mqttProxy.pending[*jsonrpc_res.ID]; pending != nil && jsonrpc_res.IsResponse() {
			mqttProxy.mutex.Lock()
			delete(mqttProxy.pending, *jsonrpc_res.ID)
			mqttProxy.mutex.Unlock()
			pending <- msg.Payload()
		}

	})

	if err := token.Error(); err != nil {
		return mqttProxy, err
	}
	return mqttProxy, nil
}

func (m *MqttProxy) Notify(id string, req []byte) error {
	token := m.Client.Publish(fmt.Sprintf(rpc_topic, id), 1, false, req)
	if err := token.Error(); err != nil {
		return err
	}
	return nil
}

func (m *MqttProxy) AsyncRpc(id string, req []byte, ch_res chan []byte) error {
	jsonrpc_req := jsonrpc2.WireRequest{}
	err := json.Unmarshal(req, &jsonrpc_req)
	if err != nil {
		return err
	}

	if jsonrpc_req.ID == nil {
		return errors.New("jsonrpc Not find ID")
	}

	m.mutex.Lock()
	m.pending[*jsonrpc_req.ID] = ch_res
	m.mutex.Unlock()

	token := m.Client.Publish(fmt.Sprintf(rpc_topic, id), 1, false, req)
	if err := token.Error(); err != nil {
		return err
	}
	return nil
}

func (m *MqttProxy) SyncRpc(id string, req []byte) ([]byte, error) {
	return m.SyncRpcWait(id, 1*time.Hour, req)
}

func (m *MqttProxy) SyncRpcWait(id string, timeout time.Duration, req []byte) ([]byte, error) {
	ch := make(chan []byte)
	err := m.AsyncRpc(id, req, ch)
	if err != nil {
		return []byte{}, err
	}
	return <-ch, nil
}
