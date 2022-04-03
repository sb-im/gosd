package cmd

import (
	"time"

	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, databaseCmd)
}

var databaseCmd = &cli.Command{
	Name:  "database",
	Usage: "Database management utility",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "local",
			Value: false,
			Usage: "Disable remote database management",
		},
	},
	Before: func(c *cli.Context) error {
		if c.Bool("local") {
			go Daemon()
			time.Sleep(time.Second)
		}
		return nil
	},
}
