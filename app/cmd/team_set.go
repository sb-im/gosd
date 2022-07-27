package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	teamCmd.AddCommand(teamSetCmd)
	teamSetCmd.Flags().String("name", "", "team name")
}

var teamSetCmd = &cobra.Command{
	Use:   "set <id>",
	Short: "Update a team information",
	Args:  cobra.ExactArgs(1),
	RunE: ex(func(c *exContext) error {
		team := make(map[string]interface{})
		if name, err := c.ctx.Flags().GetString("name"); err == nil {
			team["name"] = name
		}
		return c.cnt.TeamUpdate(c.arg[0], &team)
	}),
}
