package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, userCmd)
}

var userCmd = &cli.Command{
	Name:  "user",
	Usage: "user",
}
