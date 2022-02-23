package cmd

import (
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/config"

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
	Action: func(c *cli.Context) error {
		cc := client.NewClient(config.Opts().BaseURL, config.Opts().ApiKey)
		team := make(map[string]interface{})
		if k := c.String("name"); k != "" {
			team["name"] = k
		}
		return cc.TeamUpdate(c.Args().First(), &team)
	},
}
