package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

func init() {
	teamCmd.Subcommands = append(teamCmd.Subcommands, teamLsCmd)
}

var teamLsCmd = &cli.Command{
	Name:  "ls",
	Usage: "Ls all team",
	Action: ex(func(c *exContext) error {
		teams, _ := c.cnt.TeamIndex()
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName")

		for _, i := range teams {
			fmt.Fprintf(w, "%d\t%s\n",
				i.ID,
				i.Name,
			)
		}

		return w.Flush()
	}),
}
