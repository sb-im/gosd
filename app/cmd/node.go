package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, nodeCmd)
}

var nodeCmd = &cli.Command{
	Name:  "node",
	Usage: "Nodes management utility",
}
