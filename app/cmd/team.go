package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(teamCmd)
}

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Teams management utility",
}
