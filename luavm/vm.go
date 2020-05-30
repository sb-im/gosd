package luavm

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"sb.im/gosd/jsonrpc2mqtt"
	"sb.im/gosd/state"

	jsonrpc2 "github.com/sb-im/jsonrpc2"
	"github.com/yuin/gopher-lua"
)

var sequence uint64
var mu sync.Mutex

func getSequence() string {
	mu.Lock()
	id := strconv.FormatUint(sequence, 10)
	sequence++
	mu.Unlock()
	return id
}

func Run(s *state.State, path string) {
	l := lua.NewState()
	defer l.Close()

	regService(s, l)

	//err := l.DoString(script)
	err := l.DoFile(path)
	if err != nil {
		log.Println(err)
	}

	// 执行具体的lua脚本
	err = l.CallByParam(lua.P{
		Fn:      l.GetGlobal("info"), // 获取info函数引用
		NRet:    1,                   // 指定返回值数量
		Protect: true,                // 如果出现异常，是panic还是返回err
	}, lua.LNumber(1)) // 传递输入参数n=1
	if err != nil {
		panic(err)
	}
	// 获取返回结果
	ret := l.Get(-1)
	// 从堆栈中删除返回值
	l.Pop(1)
	// 打印返回结果
	fmt.Println(ret)
}

func regService(s *state.State, l *lua.LState) {
	mqttProxy, _ := jsonrpc2mqtt.OpenMqttProxy(s.Mqtt)
	service := LService{
		Client:    s,
		MqttProxy: mqttProxy,
	}

	l.SetGlobal("async_rpc_call", l.NewFunction(service.async_rpc))
	l.SetGlobal("rpc_call", l.NewFunction(service.rpc))
	l.SetGlobal("call_service", l.NewFunction(callService))
	l.SetGlobal("filePoolService", lua.LString("FilePoolService"))

	l.SetGlobal("plan_id", lua.LString("1"))
	l.SetGlobal("plan_log_id", lua.LString("2"))
	l.SetGlobal("node_id", lua.LString("3"))
}

type LService struct {
	Client    *state.State
	MqttProxy *jsonrpc2mqtt.MqttProxy
}

func (m *LService) async_rpc(L *lua.LState) int {
	bit13_timestamp := string([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))[:13])
	rpc_id := "gosd.0-" + bit13_timestamp + "-" + getSequence()
	fmt.Println(rpc_id)

	jsonrpc_req := jsonrpc2.WireRequest{
		Method: L.ToString(2),
		ID: &jsonrpc2.ID{
			Name: rpc_id,
		},
	}

	ch_L := L.ToChannel(3)
	ch := make(chan []byte)
	r, _ := json.Marshal(jsonrpc_req)
	err := m.MqttProxy.AsyncRpc(L.ToString(1), r, ch)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		ch_L <- lua.LString(<-ch)
	}()
	return 1
}

func (m *LService) rpc(L *lua.LState) int {
	bit13_timestamp := string([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))[:13])
	rpc_id := "gosd.0-" + bit13_timestamp + "-" + getSequence()
	fmt.Println(rpc_id)

	jsonrpc_req := jsonrpc2.WireRequest{
		Method: L.ToString(2),
		ID: &jsonrpc2.ID{
			Name: rpc_id,
		},
	}

	r, _ := json.Marshal(jsonrpc_req)
	r, err := m.MqttProxy.SyncRpc(L.ToString(1), r)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(r))
	L.Push(lua.LString(string(r)))
	return 1
}

func callService(L *lua.LState) int {
	// 根据编号获取传入参数(从1开始)
	service := L.ToString(1)
	param := L.ToTable(3)
	param.ForEach(func(key, value lua.LValue) {
		fmt.Println(key.String())
		fmt.Println(value.String())
	})

	// 注册一个table类型,设置返回参数
	t := L.NewTable()
	t.RawSet(lua.LString("msg"), lua.LString("success"))
	t.RawSet(lua.LString("data"), lua.LString(service))

	// 将返货结果堆栈
	L.Push(t)
	return 1
}
