package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, teamCmd)
}

var teamCmd = &cli.Command{
	Name:  "team",
	Usage: "Teams management utility",
}
