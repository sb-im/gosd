package cmd

import (
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/config"

	"github.com/urfave/cli/v2"
)

type exContext struct {
	ctx *cli.Context
	cnt *client.Client
}

func ex(fn func(c *exContext) error) cli.ActionFunc {
	return func(c *cli.Context) error {
		return fn(&exContext{
			ctx: c,
			cnt: client.NewClient(config.Opts().BaseURL, config.Opts().ApiKey),
		})
	}
}
