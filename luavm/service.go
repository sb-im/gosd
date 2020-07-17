package luavm

import (
	"strconv"

	"sb.im/gosd/state"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type LService struct {
	Task   *Task
	State  *state.State
	NodeID string
}

// GetMsg(id, msg string) (data tables{}, error string)
func (s *LService) GetMsg(L *lua.LState) int {
	raw, err := s.State.NodeGet(L.ToString(1), L.ToString(2))
	if err != nil {
		L.Push(&lua.LTable{})
		L.Push(lua.LString(err.Error()))
		return 2
	}

	value, err := luajson.Decode(L, raw)
	if err != nil {
		L.Push(&lua.LTable{})
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(value)
	L.Push(lua.LString(""))
	return 2
}

// GetID(str string) (id string)
func (s *LService) GetID(L *lua.LState) int {
	if n := s.State.Node[s.NodeID]; n != nil && L.ToString(1) == "link_id" {
		L.Push(lua.LString(strconv.Itoa(n.Status.GetID(""))))
		return 1
	}

	L.Push(lua.LString("0"))
	return 1
}

// GetStatus() (data tables{}, error string)
func (s *LService) GetStatus(L *lua.LState) int {
	raw := []byte{}
	if data := s.State.Node[s.NodeID]; data != nil {
		raw = data.Status.Raw
	}

	value, err := luajson.Decode(L, raw)
	if err != nil {
		L.Push(&lua.LTable{})
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(value)
	L.Push(lua.LString(""))
	return 2
}
