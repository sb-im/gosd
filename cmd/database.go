package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(databaseCmd)
}

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database management utility",
	Long:  `Database management utility.`,
	Args:  cobra.NoArgs,
}
