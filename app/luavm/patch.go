package luavm

import (
	"context"
	"fmt"

	lua "github.com/yuin/gopher-lua"
	"sb.im/gosd/app/logger"
)

// This patch hook lua base function

// Fork From: https://github.com/yuin/gopher-lua/blob/af7e27f2568ed807e8e55afc83cf8dcd88c3197f/baselib.go#L283-L293
func patchBasePrint(ctx context.Context) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		top := L.GetTop()
		str := ""
		for i := 1; i <= top; i++ {
			str += fmt.Sprint(L.ToStringMeta(L.Get(i)).String())
			if i != top {
				str += fmt.Sprint("\t")
			}
		}
		logger.WithContext(ctx).Println(str)
		return 0
	}
}
