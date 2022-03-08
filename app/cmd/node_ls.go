package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

func init() {
	nodeCmd.Subcommands = append(nodeCmd.Subcommands, nodeLsCmd)
}

var nodeLsCmd = &cli.Command{
	Name:  "ls",
	Usage: "Ls all node",
	Action: ex(func(c *exContext) error {
		nodes, err := c.cnt.NodeIndex()
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName")

		for _, i := range nodes {
			fmt.Fprintf(w, "%d\t%s\n",
				i.ID,
				i.Name,
			)
		}

		return w.Flush()
	}),
}
