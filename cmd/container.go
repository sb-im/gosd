package cmd

import (
	"sb.im/gosd/cli"
	"sb.im/gosd/config"
	"sb.im/gosd/database"
	"sb.im/gosd/model"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(containerCmd)
}

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "container mode",
	Args:  cobra.NoArgs,
	Run: ex(func(cmd *cobra.Command, args []string, d exData) {
		database.Migrate(d.store.Database())
		user, err := d.store.UserByID(1)
		if err != nil {
			panic(err)
		}

		parse := config.NewParser()
		opts, err := parse.ParseEnvironmentVariables()
		if err != nil {
			panic(err)
		}

		// not found user == nil
		if user == nil {
			// Create
			user = model.NewUser()
			user.Username = opts.AdminUsername
			user.Password = opts.AdminPassword
			// username == group
			group := model.NewGroup()
			group.Name = opts.AdminUsername
			if err := d.store.CreateGroup(group); err != nil {
				checkErr(err)
			}

			user.Group = group

			if err := d.store.CreateUser(user); err != nil {
				checkErr(err)
			}
			printUsers([]*model.User{user})
		} else {
			user.Username = opts.AdminUsername
			user.Password = opts.AdminPassword
			// Update
			err = d.store.UpdateUser(user)
			checkErr(err)
			printUsers([]*model.User{user})
		}
		if flagDebugMode || opts.HasDebugMode() {
			log.SetLevel(log.DebugLevel)
		}

		cli.StartDaemon(d.store, opts)
	}, exConfig{}),
}
