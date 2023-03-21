package cmd

import (
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/model"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func init() {
	databaseCmd.AddCommand(databaseMigrateCmd)
}

var databaseMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate",
	Args:  cobra.ExactArgs(0),
	RunE: func(c *cobra.Command, args []string) error {
		orm, err := DatabaseOrm(config.Opts())
		if err != nil {
			return err
		}
		DatabaseMigrate(orm)
		log.Warn("=== Database Migrate Done ===")
		return nil
	},
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
