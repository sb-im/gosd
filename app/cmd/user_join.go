package cmd

import (
	"sb.im/gosd/app/client"
	"sb.im/gosd/app/config"

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
	Action: func(c *cli.Context) error {
		cc := client.NewClient(config.Opts().BaseURL, config.Opts().ApiKey)
		return cc.UserAddTeam(c.Args().First(), c.String("team"))
	},
}
