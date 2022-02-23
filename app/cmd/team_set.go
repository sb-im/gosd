package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	teamCmd.Subcommands = append(teamCmd.Subcommands, teamSetCmd)
}

var teamSetCmd = &cli.Command{
	Name:  "set",
	Usage: "Update a team information",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
	},
	ArgsUsage: "<id>",
	Action: ex(func(c *exContext) error {
		team := make(map[string]interface{})
		if k := c.ctx.String("name"); k != "" {
			team["name"] = k
		}
		return c.cnt.TeamUpdate(c.ctx.Args().First(), &team)
	}),
}
