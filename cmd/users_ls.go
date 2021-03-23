package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	usersCmd.AddCommand(usersLsCmd)
}

var usersLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all users.",
	Args:  cobra.NoArgs,
	Run: ex(func(cmd *cobra.Command, args []string, d exData) {
		users, err := d.store.Users()
		if err != nil {}
		printUsers(users)
	}, exConfig{}),
}
