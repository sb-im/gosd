package cli

import (
	"fmt"
	"os"

	"sb.im/gosd/model"
	"sb.im/gosd/storage"
)

func createAdmin(store *storage.Storage) {
	user := model.NewUser()
	user.Username = os.Getenv("ADMIN_USERNAME")
	user.Password = os.Getenv("ADMIN_PASSWORD")

	group := model.NewGroup()
	group.Name = os.Getenv("ADMIN_GROUP")

	if user.Username == "" || user.Password == "" {
		user.Username, user.Password, group.Name = askCredentials()
	}

	if err := user.ValidateUserCreation(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if store.UserExists(user.Username) {
		fmt.Printf(`User %q already exists, skipping creation`, user.Username)
		return
	}

	if err := store.CreateGroup(group); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	user.Group = group

	if err := store.CreateUser(user); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
