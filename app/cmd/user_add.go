package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(userAddCmd)
	userAddCmd.Flags().Uint("team", 0, "Team ID")
	userAddCmd.Flags().StringP("username", "u", "", "Username")
	userAddCmd.Flags().StringP("password", "p", "", "Password")
	userAddCmd.Flags().String("language", "", "Language")
	userAddCmd.Flags().String("timezone", "", "Timezone")
}

var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new user",
	RunE: ex(func(c *exContext) error {
		user := &map[string]interface{}{
			"team_id":  mustGetUint(c.ctx.Flags(), "team"),
			"username": mustGetString(c.ctx.Flags(), "username"),
			"password": mustGetString(c.ctx.Flags(), "password"),
			"language": mustGetString(c.ctx.Flags(), "language"),
			"timezone": mustGetString(c.ctx.Flags(), "timezone"),
		}

		return c.cnt.UserCreate(user)
	}),
}
