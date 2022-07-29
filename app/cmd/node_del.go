package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	nodeCmd.AddCommand(nodeDelCmd)
}

var nodeDelCmd = &cobra.Command{
	Use:   "del <id|uuid>",
	Short: "Delete a node",
	RunE: ex(func(c *exContext) error {
		return c.cnt.NodeDestroy(c.arg[0])
	}),
}
