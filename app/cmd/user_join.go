package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	userCmd.Subcommands = append(userCmd.Subcommands, userJoinCmd)
}

var userJoinCmd = &cli.Command{
	Name:  "join",
	Usage: "Join user to team",
	Flags: []cli.Flag{
		&cli.UintFlag{Name: "team"},
	},
	ArgsUsage: "<user id>",
	Action: ex(func(c *exContext) error {
		return c.cnt.UserAddTeam(c.ctx.Args().First(), c.ctx.String("team"))
	}),
}
