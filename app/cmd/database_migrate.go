package cmd

import (
	"sb.im/gosd/app/model"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func init() {
	databaseCmd.Subcommands = append(databaseCmd.Subcommands, &cli.Command{
		Name:  "migrate",
		Usage: "migrate",
		Action: func(c *cli.Context) error {
			orm, err := DatabaseOrm()
			if err != nil {
				return err
			}
			DatabaseMigrate(orm)
			log.Warn("=== Database Migrate Done ===")
			return nil
		},
	})
}

func DatabaseMigrate(orm *gorm.DB) {
	orm.AutoMigrate(&model.Team{})
	orm.AutoMigrate(&model.User{})
	orm.AutoMigrate(&model.Session{})
	orm.AutoMigrate(&model.UserTeam{})

	orm.AutoMigrate(&model.Schedule{})
	orm.AutoMigrate(&model.Node{})
	orm.AutoMigrate(&model.Task{})
	orm.AutoMigrate(&model.Blob{})
	orm.AutoMigrate(&model.Job{})
	orm.AutoMigrate(&model.Profile{})
}
