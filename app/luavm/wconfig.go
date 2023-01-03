package luavm

const (
	defaultFileLua = "default.lua"
	defaultFileKey = "lua"
)

// Order is important
var (
	libs = []string{
		"lib_task.lua",
		"lib_node.lua",
		"lib_geo.lua",
		"lib_log.lua",
		"lib_main.lua",
	}
)

type Config struct {
	Instance string
	BaseURL  string
}
