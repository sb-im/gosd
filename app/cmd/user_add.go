package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	userCmd.Subcommands = append(userCmd.Subcommands, userAddCmd)
}

var userAddCmd = &cli.Command{
	Name:  "add",
	Usage: "Create a new user",
	Flags: []cli.Flag{
		&cli.UintFlag{Name: "team"},
		&cli.StringFlag{Name: "username", Aliases: []string{"u"}},
		&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
		&cli.StringFlag{Name: "language"},
		&cli.StringFlag{Name: "timezone"},
	},
	Action: ex(func(c *exContext) error {
		user := &map[string]interface{}{
			"team_id":  c.ctx.Uint("team"),
			"username": c.ctx.String("username"),
			"password": c.ctx.String("password"),
			"language": c.ctx.String("language"),
			"timezone": c.ctx.String("timezone"),
		}

		return c.cnt.UserCreate(user)
	}),
}
