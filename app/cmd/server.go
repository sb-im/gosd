package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sb.im/gosd/app/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Up Daemon",
	RunE: ex(func(c *exContext) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		go Daemon(ctx)

		cfg := config.Parse()
		if cfg.DemoMode {
			time.Sleep(3 * time.Second)
			c.cnt.NodeSync(mustGetUint(c.ctx.Flags(), "team"), "data")
		}

		select {
		case <-ctx.Done():
		}
		log.Warn("=== Safe Down Server ===")
		log.Warn("=== END ===")
		return nil
	}),
}
