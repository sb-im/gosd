package cmd

import (
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/daemon"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		daemon.DatabaseMigrate(orm)
		log.Warn("=== Database Migrate Done ===")
		return nil
	},
}
