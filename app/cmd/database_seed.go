package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	databaseCmd.Subcommands = append(databaseCmd.Subcommands, &cli.Command{
		Name:  "seed",
		Usage: "seed",
		Action: ex(func(c *exContext) error {
			c.cnt.DatabaseSeed()
			return nil
		}),
	})
}
