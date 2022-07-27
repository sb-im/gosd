package cmd

import (
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/config"

	"github.com/spf13/cobra"
)

type exContext struct {
	ctx *cobra.Command
	cnt *client.Client
	arg []string
}

func ex(fn func(c *exContext) error) func(cmd *cobra.Command, args []string) error {
	return func(c *cobra.Command, args []string) error {
		return fn(&exContext{
			ctx: c,
			cnt: client.NewClient(config.Opts().BaseURL, config.Opts().ApiKey),
			arg: args,
		})
	}
}
