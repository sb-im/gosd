package cmd

import (
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/config"

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
	Action: func(c *cli.Context) error {
		cc := client.NewClient(config.Opts().BaseURL)
		user := make(map[string]interface{})
		if k := c.Uint("team"); k != 0 {
			user["team"] = k
		}
		if k := c.String("username"); k != "" {
			user["username"] = k
		}
		if k := c.String("password"); k != "" {
			user["password"] = k
		}
		if k := c.String("language"); k != "" {
			user["language"] = k
		}
		if k := c.String("timezone"); k != "" {
			user["timezone"] = k
		}
		return cc.UserUpdate(c.Args().First(), &user)
	},
}
