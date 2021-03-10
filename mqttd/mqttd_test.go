package mqttd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"sb.im/gosd/state"
)

func cmdRun(str string) {
	fmt.Println("EXEC:", str)
	if out, err := helpCmdRun(str); err != nil {
		fmt.Printf("%s", out)
		panic(err)
	}
}

func helpCmdRun(str string) ([]byte, error) {
	cmdArr := strings.Split(str, " ")
	return exec.Command(cmdArr[0], cmdArr[1:]...).CombinedOutput()
}

func helpGetMqttAddr() string {
	mqttAddr := "mqtt://localhost:1883"
	if addr := os.Getenv("MQTT"); addr != "" {
		mqttAddr = fmt.Sprintf("mqtt://%s", addr)
	}
	return mqttAddr
}

func helpGetRedisAddr() string {
	redisAddr := "redis://localhost:6379/0"
	if addr := os.Getenv("REDIS_URL"); addr != "" {
		redisAddr = addr
	}
	return redisAddr
}

func TestMqttd(t *testing.T) {
	id := "000"
	mqttAddr := helpGetMqttAddr()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chI := make(chan MqttRpc)
	chO := make(chan MqttRpc)

	store := state.NewState(helpGetRedisAddr())
	mqttd := NewMqttd(mqttAddr, store, chI, chO)
	go mqttd.Run(ctx)

	// Wait for mqttd to start and subscribe successfully
	time.Sleep(1 * time.Second)

	// Msg
	rawMsg := `{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}`
	cmdRun("mosquitto_pub -L " + mqttAddr + "/nodes/" + id + "/msg/weather -m " + rawMsg)

	if msg, err := store.GetNodeMsg(id, "weather"); err != nil {
		t.Error(err)
	} else if string(msg) != rawMsg {
		t.Errorf("%s\n", msg)
	}
}

func TestMqttdRpc(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	id := "000"
	mqttRpcRecv, mqttRpcSend := "nodes/%s/rpc/recv", "nodes/%s/rpc/send"
	mqttAddr := helpGetMqttAddr()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chI := make(chan MqttRpc)
	chO := make(chan MqttRpc)

	store := state.NewState(helpGetRedisAddr())
	mqttd := NewMqttd(mqttAddr, store, chI, chO)
	go mqttd.Run(ctx)
	//time.Sleep(3 * time.Second)

	// RPC

	// Send
	rawRpcSend := `{"jsonrpc":"2.0","method":"test","id":"test.%d"}`

	// Sub
	sub := exec.CommandContext(ctx, "mosquitto_sub", "-L", mqttAddr+"/"+fmt.Sprintf(mqttRpcSend, id))
	stdout, err := sub.StdoutPipe()
	if err != nil {
		t.Error(err)
	}
	if err := sub.Start(); err != nil {
		t.Error(err)
	}
	// Wait sub topic on broker
	time.Sleep(10 * time.Millisecond)

	// Pub
	req := fmt.Sprintf(rawRpcSend, 233)
	chI <- MqttRpc{
		ID:      id,
		Payload: []byte(req),
	}

	// ValidateSend
	reader := bufio.NewReader(stdout)
	for i := 0; i < 1; i++ {
		raw, err := reader.ReadString('\n')
		if err != nil {
			t.Error(err)
		}
		if res := strings.TrimSuffix(raw, "\n"); res != req {
			t.Errorf("Recv is: %s, Should: %s", res, req)
		}
	}

	// Recv
	rawRpcRecv := `{"jsonrpc":"2.0","result":"ok","id":"test.%d"}`
	res2 := fmt.Sprintf(rawRpcRecv, 234)
	cmdRun("mosquitto_pub -L " + mqttAddr + "/" + fmt.Sprintf(mqttRpcRecv, id) + " -m " + res2)

	// ValidateRecv
	p := <-chO
	//fmt.Printf("%s\n", p)
	if p.ID != id {
		t.Error("id is: ", p.ID)
	}

	if string(p.Payload) != res2 {
		t.Error(res2)
	}

	time.Sleep(1 * time.Second)
}
