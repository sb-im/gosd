package luavm

import "os"

type Config struct {
	Instance string
	LuaFile  string
	LuaTask  string
	Timeout  string
	BaseURL  string
}

func DefaultConfig() Config {
	return Config{
		Instance: "gosd.0",
		LuaFile:  "test_min.lua",
		LuaTask:  "lua",
		Timeout:  "2h",
		BaseURL:  os.Getenv("BASE_URL"),
	}
}
