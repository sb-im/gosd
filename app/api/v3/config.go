package v3

import (
	"sb.im/gosd/app/model"
)

var (
	DefaultConfig = &Config{
		Auth: ConfigAuth{
			SuperAdmin: map[uint]*model.User{
				1: {
					Username: "superadmin",
					Password: "superadmin",
				},
			},
		},
		StoragePath: "data/storage/",
	}
)

type Config struct {
	Auth        ConfigAuth
	StoragePath string
}

type ConfigAuth struct {
	SuperAdmin map[uint]*model.User
	JWTSecret  string
}
