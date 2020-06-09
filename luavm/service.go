package luavm

import (
	"sb.im/gosd/state"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type LService struct {
	State *state.State
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

