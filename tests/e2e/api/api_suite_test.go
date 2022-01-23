package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
	t.Setenv("LUA_FILE", "app/luavm/lua/test_rpc.lua")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
