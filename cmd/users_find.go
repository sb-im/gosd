package cmd

import (
	users "sb.im/gosd/model"

	"github.com/spf13/cobra"
)

func init() {
	usersCmd.AddCommand(usersFindCmd)
}

var usersFindCmd = &cobra.Command{
	Use:   "find <id|username>",
	Short: "Find a user by username or id",
	Long:  `Find a user by username or id. If no flag is set, all users will be printed.`,
	Args:  cobra.ExactArgs(1),
	Run: findUsers,
}

var findUsers = ex(func(cmd *cobra.Command, args []string, d exData) {
	var (
		list []*users.User
		user *users.User
		err  error
	)

	if len(args) == 1 {
		username, id := parseUsernameOrID(args[0])
		if username != "" {
			user, err = d.store.UserByUsername(username)
		} else {
			user, err = d.store.UserByID(id)
		}

		list = []*users.User{user}
	} else {
		list, err = d.store.Users()
	}

	// not found user == nil
	if user == nil { return }

	checkErr(err)
	printUsers(list)
}, exConfig{})
