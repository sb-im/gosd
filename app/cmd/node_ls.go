package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
		type point struct {
			NodeID int             `json:"node_id"`
			Params json.RawMessage `json:"params"`
			Type   string          `json:"type"`
		}
		nodes, err := c.cnt.NodeIndex()
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tPoints")

		for _, i := range nodes {
			fmt.Fprintf(w, "%d\t%s\t%s\n",
				i.ID,
				i.Name,
				func() string {
					var points []point
					if err := json.Unmarshal(i.Points, &points); err != nil {
						panic(err)
					}
					arr := make([]string, len(points))
					for i, t := range points {
						arr[i] = t.Type
					}
					return strings.Join(arr, ",")
				}(),
			)
		}

		return w.Flush()
	}),
}
