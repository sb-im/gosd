package cmd

import (
	"fmt"

	"sb.im/gosd/model"

	"github.com/spf13/cobra"
)

func init() {
	usersCmd.AddCommand(usersAddCmd)
}

var usersAddCmd = &cobra.Command{
	Use:   "add <username> <password>",
	Short: "Create a new user",
	Long:  `Create a new user and add it to the database.`,
	Args:  cobra.ExactArgs(2), //nolint:gomnd
	Run: ex(func(cmd *cobra.Command, args []string, d exData) {
		user := model.NewUser()
		user.Username = args[0]
		user.Password = args[1]

		if d.store.UserExists(user.Username) {
			fmt.Printf(`User %q already exists, skipping creation`, user.Username)
			return
		}

		// username == group
		group := model.NewGroup()
		group.Name = args[0]
		if err := d.store.CreateGroup(group); err != nil {
			checkErr(err)
		}

		user.Group = group

		if err := d.store.CreateUser(user); err != nil {
			checkErr(err)
		}
		printUsers([]*model.User{user})
	}, exConfig{}),
}
