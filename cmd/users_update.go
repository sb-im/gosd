package cmd

import (
	users "sb.im/gosd/model"

	"github.com/spf13/cobra"
)

func init() {
	usersCmd.AddCommand(usersUpdateCmd)

	usersUpdateCmd.Flags().StringP("password", "p", "", "new password")
	usersUpdateCmd.Flags().StringP("username", "u", "", "new username")
}

var usersUpdateCmd = &cobra.Command{
	Use:   "update <id|username>",
	Short: "Updates an existing user",
	Long: `Updates an existing user. Set the flags for the
options you want to change.`,
	Args: cobra.ExactArgs(1),
	Run: ex(func(cmd *cobra.Command, args []string, d exData) {
		username, id := parseUsernameOrID(args[0])
		flags := cmd.Flags()
		password := mustGetString(flags, "password")
		newUsername := mustGetString(flags, "username")

		var (
			err  error
			user *users.User
		)

		if username != "" {
			user, _ = d.store.UserByUsername(username)
		} else {
			user, _ = d.store.UserByID(id)
		}
		checkErr(err)
		// not found user == nil
		if user == nil { return }

		if newUsername != "" {
			user.Username = newUsername
		}

		if password != "" {
			user.Password = password
		}

		err = d.store.UpdateUser(user)
		checkErr(err)
		printUsers([]*users.User{user})
	}, exConfig{}),
}
