package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	nodeCmd.AddCommand(nodeSyncCmd)
	nodeSyncCmd.Flags().Uint("team", 0, "Team ID")
}

var nodeSyncCmd = &cobra.Command{
	Use:   "sync <path>",
	Short: "Create Or Update batch nodes",
	Args:  cobra.ExactArgs(1),
	RunE: ex(func(c *exContext) error {
		return c.cnt.NodeSync(mustGetUint(c.ctx.Flags(), "team"), c.arg[0])
	}),
}
