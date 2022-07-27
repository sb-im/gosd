package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nodeCmd)
}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Nodes management utility",
}
