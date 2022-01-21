package cmd

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/model"
)

func init() {
	databaseCmd.Subcommands = append(databaseCmd.Subcommands, &cli.Command{
		Name:  "migrate",
		Usage: "migrate",
		Action: func(c *cli.Context) error {
			databaseMigrate()
			return nil
		},
	})
}

func databaseMigrate() {
	cfg := config.DefaultConfig()
	if err := env.Parse(cfg); err != nil {
		log.Errorf("%+v\n", err)
	}

	orm, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	orm.AutoMigrate(&model.Team{})
	orm.AutoMigrate(&model.User{})
	orm.AutoMigrate(&model.Session{})
	orm.AutoMigrate(&model.UserTeam{})

	orm.AutoMigrate(&model.Schedule{})
	orm.AutoMigrate(&model.Task{})
	orm.AutoMigrate(&model.Blob{})
	orm.AutoMigrate(&model.Job{})

	log.Warn("=== Database Migrate Done ===")
}
