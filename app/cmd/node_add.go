package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	nodeCmd.Subcommands = append(nodeCmd.Subcommands, nodeAddCmd)
}

var nodeAddCmd = &cli.Command{
	Name:  "add",
	Usage: "Create a new node",
	Flags: []cli.Flag{
		&cli.UintFlag{Name: "team"},
		&cli.StringFlag{Name: "name"},
	},
	Action: ex(func(c *exContext) error {
		return c.cnt.NodeCreate(&map[string]interface{}{
			"team_id": c.ctx.Uint("team"),
			"name":    c.ctx.String("name"),
		})
	}),
}
