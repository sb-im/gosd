package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(userSetCmd)
	userSetCmd.Flags().Uint("team", 0, "Team ID")
	userSetCmd.Flags().StringP("username", "u", "", "Username")
	userSetCmd.Flags().StringP("password", "p", "", "Password")
	userSetCmd.Flags().String("language", "", "Language")
	userSetCmd.Flags().String("timezone", "", "Timezone")
}

var userSetCmd = &cobra.Command{
	Use:   "set <id>",
	Short: "Update a user information",
	Args:  cobra.ExactArgs(1),
	RunE: ex(func(c *exContext) error {
		user := make(map[string]interface{})
		if k, _ := c.ctx.Flags().GetUint("team"); k != 0 {
			user["team"] = k
		}
		if k, _ := c.ctx.Flags().GetString("username"); k != "" {
			user["username"] = k
		}
		if k, _ := c.ctx.Flags().GetString("password"); k != "" {
			user["password"] = k
		}
		if k, _ := c.ctx.Flags().GetString("language"); k != "" {
			user["language"] = k
		}
		if k, _ := c.ctx.Flags().GetString("timezone"); k != "" {
			user["timezone"] = k
		}
		return c.cnt.UserUpdate(c.arg[0], &user)
	}),
}
