package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	teamCmd.Subcommands = append(teamCmd.Subcommands, teamAddCmd)
}

var teamAddCmd = &cli.Command{
	Name:  "add",
	Usage: "Create a new team",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
	},
	Action: ex(func(c *exContext) error {
		team := &map[string]interface{}{
			"name": c.ctx.String("name"),
		}

		return c.cnt.TeamCreate(team)
	}),
}
