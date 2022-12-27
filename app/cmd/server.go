package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"sb.im/gosd/app/config"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Up Daemon",
	RunE: ex(func(c *exContext) error {
		go Daemon()

		cfg := config.Parse()
		if cfg.DemoMode {
			time.Sleep(3 * time.Second)
			c.cnt.NodeSync(mustGetUint(c.ctx.Flags(), "team"), "data")
		}

		select {}
	}),
}
