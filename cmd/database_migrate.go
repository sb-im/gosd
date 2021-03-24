package cmd

import (
	"sb.im/gosd/database"

	"github.com/spf13/cobra"
)

func init() {
	databaseCmd.AddCommand(databaseMigrateCmd)
}

var databaseMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "database migrate",
	Args:  cobra.NoArgs,
	Run: ex(func(cmd *cobra.Command, args []string, d exData) {
		database.Migrate(d.store.Database())
	}, exConfig{}),
}
