package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(userJoinCmd)
	userJoinCmd.Flags().String("team", "", "Team ID")
}

var userJoinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join user to team",
	Args:  cobra.ExactArgs(1),
	RunE: ex(func(c *exContext) error {
		return c.cnt.UserAddTeam(c.arg[0], mustGetString(c.ctx.Flags(), "team"))
	}),
}
