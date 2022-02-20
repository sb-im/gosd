package cmd

import (
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/config"

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
	Action: func(c *cli.Context) error {
		cc := client.NewClient(config.Opts().BaseURL)
		team := &map[string]interface{}{
			"name": c.String("name"),
		}

		return cc.TeamCreate(team)
	},
}
