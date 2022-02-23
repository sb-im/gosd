package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	userCmd.Subcommands = append(userCmd.Subcommands, userSetCmd)
}

var userSetCmd = &cli.Command{
	Name:  "set",
	Usage: "Update a user information",
	Flags: []cli.Flag{
		&cli.UintFlag{Name: "team"},
		&cli.StringFlag{Name: "username", Aliases: []string{"u"}},
		&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
		&cli.StringFlag{Name: "language"},
		&cli.StringFlag{Name: "timezone"},
	},
	ArgsUsage: "<id>",
	Action: ex(func(c *exContext) error {
		user := make(map[string]interface{})
		if k := c.ctx.Uint("team"); k != 0 {
			user["team"] = k
		}
		if k := c.ctx.String("username"); k != "" {
			user["username"] = k
		}
		if k := c.ctx.String("password"); k != "" {
			user["password"] = k
		}
		if k := c.ctx.String("language"); k != "" {
			user["language"] = k
		}
		if k := c.ctx.String("timezone"); k != "" {
			user["timezone"] = k
		}
		return c.cnt.UserUpdate(c.ctx.Args().First(), &user)
	}),
}
