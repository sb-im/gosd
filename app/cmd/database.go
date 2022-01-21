package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, databaseCmd)
}

var databaseCmd = &cli.Command{
	Name:  "database",
	Usage: "database",
}
