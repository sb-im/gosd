package cmd

import (
	"os"

	"sb.im/gosd/version"

	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var (
	app = &cli.App{
		EnableBashCompletion: true,

		Name:    "gosd",
		Version: version.Version + " " + version.Date,
		Usage:   "SuperDock Cloud Service",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "Enabled Debug mode",
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("debug") {
				log.SetReportCaller(true)
				log.SetLevel(log.DebugLevel)
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			Daemon()
			return nil
		},
	}
)

func Execute() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
