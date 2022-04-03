package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	databaseCmd.Subcommands = append(databaseCmd.Subcommands, &cli.Command{
		Name:  "migrate",
		Usage: "migrate",
		Action: ex(func(c *exContext) error {
			c.cnt.DatabaseMigrate()
			return nil
		}),
	})
}
