package cmd

import (
	"os"

	"sb.im/gosd/version"

	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var (
	app = &cli.App{
		Name:    "gosd",
		Version: version.Version + " " + version.Date,
		Usage:   "SuperDock Cloud Service",
		Flags: []cli.Flag{
			//&cli.StringFlag{
			//	Name:    "config",
			//	Aliases: []string{"c"},
			//	Value:   "gosd.toml",
			//	Usage:   "Load configuration from `FILE`",
			//},
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
			//log.Debugln(c.Path("config"))
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
