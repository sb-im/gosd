package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"sb.im/gosd/config"
	"sb.im/gosd/database"
	"sb.im/gosd/luavm"
	"sb.im/gosd/mqttd"
	"sb.im/gosd/rpc2mqtt"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"
)

func helpGenerateMqttConfig(name string, config []byte) {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	if _, err := file.Write(config); err != nil {
		panic(err)
	}
}

func CmdRun(str string) ([]byte, error) {
	cmdArr := strings.Split(str, " ")
	return exec.Command(cmdArr[0], cmdArr[1:]...).CombinedOutput()
}

func startupNcp(id string) {
	mqttRpcRecv, mqttRpcSend := "nodes/%s/rpc/recv", "nodes/%s/rpc/send"
	mqttAddr := "mqtt://localhost:1883"
	if addr := os.Getenv("MQTT"); addr != "" {
		// addr "localhost:1883"
		mqttAddr = fmt.Sprintf("mqtt://%s", addr)
	}

	mqttdConfigPath := "/tmp/test_mqttd.yml"
	var mqttdConfig = `
mqttd:
  id: ` + id + `
  static:
    link_id: 1
    lat: "22.6876423001"
    lng: "114.2248673001"
    alt: "10088.0001"
  client: "node-%s"
  status:  "nodes/%s/status"
  network: "nodes/%s/network"
  broker: ` + mqttAddr + `
  rpc :
    i: ` + mqttRpcRecv + `
    o: ` + mqttRpcSend + `
  gtran:
    prefix: "nodes/%s/msg/%s"
  trans:
    wether:
      retain: true
      qos: 0
    battery:
      retain: true
      qos: 0

ncpio:
  - type: mqtt
    params: ` + mqttdConfigPath + `
    i_rules:
      - regexp: '.*'
    o_rules:
      - regexp: '.*'
  - type: jsonrpc2
    params: "233"
    i_rules:
      - regexp: '.*'
    o_rules:
      - regexp: '.*'

`
	helpGenerateMqttConfig(mqttdConfigPath, []byte(mqttdConfig))

	go CmdRun("ncp -c " + mqttdConfigPath)

	// Wait mqttd server startup && sub topic on broker
	time.Sleep(3 * time.Millisecond)

	// Wait load mqttdConfig after delete
	os.Remove(mqttdConfigPath)
}

func TestIntegration(t *testing.T) {
	id := "0"
	parse := config.NewParser()
	opts, err := parse.ParseEnvironmentVariables()
	if err != nil {
		panic(err)
	}

	db, err := database.NewConnectionPool(
		opts.DatabaseURL(),
		opts.DatabaseMinConns(),
		opts.DatabaseMaxConns(),
	)
	store := storage.NewStorage(db)

	state := state.NewState(opts.RedisURL())

	// Start NCP
	go startupNcp(id)

	// Wait mqtt connected
	time.Sleep(3 * time.Second)

	chI := make(chan mqttd.MqttRpc, 128)
	chO := make(chan mqttd.MqttRpc, 128)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mqttd := mqttd.NewMqttd(opts.MqttURL(), state, chI, chO)
	go mqttd.Run(ctx)

	rpcServer := rpc2mqtt.NewRpc2Mqtt(chI, chO)
	go rpcServer.Run(ctx)

	worker := luavm.NewWorker(state, store, rpcServer)
	go worker.Run()

	// === start ===
	file, err := os.Open("../luavm/lua/test_rpc.lua")
	if err != nil {
		t.Error(err)
	}

	script, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	id64, err := strconv.ParseInt(id, 10, 64)
	worker.Queue <- &luavm.Task{
		NodeID: id64,
		URL:    "1/12/3/4/4",
		Script: script,
	}

	time.Sleep(3 * time.Second)
}
