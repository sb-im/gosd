package luavm

import (
	"fmt"
	"log"

	"sb.im/gosd/state"

	"github.com/yuin/gopher-lua"
)

func Run(s *state.State, path string) {
	l := lua.NewState()
	defer l.Close()

	regService(l)

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

func regService(l *lua.LState) {
	aaa := &aa{}

	l.SetGlobal("call_service", l.NewFunction(aaa.callService))
	l.SetGlobal("filePoolService", lua.LString("FilePoolService"))
}

type aa struct {
}

// lua脚本中调用的函数
func (a *aa) callService(L *lua.LState) int {
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
