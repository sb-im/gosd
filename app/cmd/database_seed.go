package cmd

import (
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/daemon"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	databaseCmd.AddCommand(databaseSeedCmd)
}

var databaseSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed",
	Args:  cobra.ExactArgs(0),
	RunE: func(c *cobra.Command, args []string) error {
		orm, err := DatabaseOrm(config.Opts())
		if err != nil {
			return err
		}
		err = daemon.DatabaseSeed(orm)
		log.Warn("=== Database Seed Done ===")
		return err
	},
}
