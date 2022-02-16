package v3

import (
	"sb.im/gosd/app/model"
)

var (
	DefaultConfig = &Config{
		Auth: ConfigAuth{
			JWTSecret: "secretkey",
			SuperAdmin: map[uint]*model.User{
				1: {
					Username: "admin",
					Password: "admin",
				},
			},
		},
		SingleUserMode:  true,
		DefaultLanguage: "en_US",
		DefaultTimezone: "Asia/Shanghai",
		StoragePath:     "data/storage/",
	}
)

type Config struct {
	SingleUserMode  bool `yaml:"SingleUserMode"`
	DefaultLanguage string
	DefaultTimezone string

	Auth        ConfigAuth
	StoragePath string
}

type ConfigAuth struct {
	SuperAdmin map[uint]*model.User
	JWTSecret  string
}
