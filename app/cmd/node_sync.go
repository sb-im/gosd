package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	nodeCmd.Subcommands = append(nodeCmd.Subcommands, nodeSyncCmd)
}

var nodeSyncCmd = &cli.Command{
	Name:  "sync",
	Usage: "Create Or Update batch nodes",
	Flags: []cli.Flag{
		&cli.UintFlag{Name: "team"},
	},
	ArgsUsage: "<path>",
	Action: ex(func(c *exContext) error {
		return c.cnt.NodeSync(c.ctx.Uint("team"), c.ctx.Args().First())
	}),
}
