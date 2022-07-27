package cmd

import (
	"sb.im/gosd/app/config"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	rootCmd.AddCommand(databaseCmd)
}

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database management utility",
}

func DatabaseOrm() (*gorm.DB, error) {
	cfg := config.Opts()

	return gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
}
