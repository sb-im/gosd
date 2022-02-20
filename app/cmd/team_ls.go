package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"sb.im/gosd/app/client"
	"sb.im/gosd/app/config"

	"github.com/urfave/cli/v2"
)

func init() {
	teamCmd.Subcommands = append(teamCmd.Subcommands, teamLsCmd)
}

var teamLsCmd = &cli.Command{
	Name:  "ls",
	Usage: "ls all team",
	Action: func(c *cli.Context) error {
		cc := client.NewClient(config.Opts().BaseURL)
		teams, _ := cc.TeamIndex()
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName")

		for _, i := range teams {
			fmt.Fprintf(w, "%d\t%s\n",
				i.ID,
				i.Name,
			)
		}

		return w.Flush()
	},
}
