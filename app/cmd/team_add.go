package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	teamCmd.AddCommand(teamAddCmd)
	teamAddCmd.Flags().String("name", "", "team name")
}

var teamAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new team",
	RunE: ex(func(c *exContext) error {
		team := &map[string]interface{}{
			"name": mustGetString(c.ctx.Flags(), "name"),
		}

		return c.cnt.TeamCreate(team)
	}),
}
