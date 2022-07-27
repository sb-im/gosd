package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	teamCmd.AddCommand(teamLsCmd)
}

var teamLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Ls all team",
	RunE: ex(func(c *exContext) error {
		teams, err := c.cnt.TeamIndex()
		if err != nil {
			return err
		}
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
