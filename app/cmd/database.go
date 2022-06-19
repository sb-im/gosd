package cmd

import (
	"sb.im/gosd/app/config"

	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
}

func DatabaseOrm() (*gorm.DB, error) {
	cfg := config.Opts()

	return gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
}
