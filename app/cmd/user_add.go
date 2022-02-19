package cmd

import (
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/config"

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
	Action: func(c *cli.Context) error {
		cc := client.NewClient(config.Opts().BaseURL)
		user := &map[string]interface{}{
			"team_id":  c.Uint("team"),
			"username": c.String("username"),
			"password": c.String("password"),
			"language": c.String("language"),
			"timezone": c.String("timezone"),
		}

		return cc.UserCreate(user)
	},
}
