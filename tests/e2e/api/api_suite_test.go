package api_test

import (
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sb.im/gosd/app/helper"
)

func TestApi(t *testing.T) {
	t.Setenv("LUA_FILE", "app/luavm/lua/test_rpc.lua")
	t.Setenv("STORAGE_URL", "file://"+path.Join(os.TempDir(), "gosd_storage_test_tmp", helper.GenNumberSecret(8)))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
