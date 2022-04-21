package help

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

func CmdRun(ctx context.Context, str string) ([]byte, error) {
	cmdArr := strings.Split(str, " ")
	return exec.CommandContext(ctx, cmdArr[0], cmdArr[1:]...).CombinedOutput()
}

func StartNcp(ctx context.Context, mqttAddr, id string) error {
	mqttRpcRecv, mqttRpcSend := "nodes/%s/rpc/recv", "nodes/%s/rpc/send"
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
      - regexp: '.*"method": ?"(__luavm_test__no_result)".*'
        invert: true
      - regexp: '.*'
    o_rules:
      - regexp: '.*'

`

	if err := os.WriteFile(mqttdConfigPath, []byte(mqttdConfig), 0666); err != nil {
		return err
	}
	defer os.Remove(mqttdConfigPath)

	_, err := CmdRun(ctx, "ncp -c "+mqttdConfigPath)
	return err
}
