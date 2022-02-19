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
		&cli.StringFlag{Name: "username", Aliases: []string{"u"}},
		&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
		&cli.StringFlag{Name: "language", Aliases: []string{"l"}},
		&cli.StringFlag{Name: "timezone", Aliases: []string{"t"}},
	},
	Action: func(c *cli.Context) error {
		cc := client.NewClient(config.Opts().BaseURL)
		user := &map[string]interface{}{
			//TODO: team
			"team_id":  1,
			"username": c.String("username"),
			"password": c.String("password"),
			"language": c.String("language"),
			"timezone": c.String("timezone"),
		}

		return cc.UserCreate(user)
	},
}
