package jsonrpc2mqtt

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsonrpc2 "github.com/sb-im/jsonrpc2"
)

func SyncMqttRpc(client mqtt.Client, id int, req []byte) ([]byte, error) {
	return SyncMqttRpcWait(client, id, 1*time.Hour, req)
}

func SyncMqttRpcWait(client mqtt.Client, id int, timeout time.Duration, req []byte) ([]byte, error) {
	topic := "nodes/" + strconv.Itoa(id) + "/rpc/"
	jsonrpc_req := jsonrpc2.WireRequest{}
	err := json.Unmarshal(req, &jsonrpc_req)
	if err != nil {
		return []byte(""), err
	}

	ch_recv := make(chan []byte)
	token := client.Subscribe(topic+"recv", 1, func(client mqtt.Client, mqtt_msg mqtt.Message) {
		jsonrpc_res := jsonrpc2.WireResponse{}
		err := json.Unmarshal(mqtt_msg.Payload(), &jsonrpc_res)
		if err != nil {
			return
		}

		if *jsonrpc_req.ID == *jsonrpc_res.ID {
			client.Unsubscribe(topic + "recv")
			ch_recv <- mqtt_msg.Payload()
		}
	})

	if token.Error() != nil {
		return []byte(""), token.Error()
	}

	token = client.Publish(topic+"send", 2, false, string(req))
	if token.Error() != nil {
		return []byte(""), token.Error()
	}

	select {
	case <-time.After(timeout):
		return []byte(""), errors.New("JSONRPC timeout")
	case result := <-ch_recv:
		return result, nil
	}
}
