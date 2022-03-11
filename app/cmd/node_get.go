package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/TylerBrock/colorjson"
	"github.com/urfave/cli/v2"
)

func init() {
	nodeCmd.Subcommands = append(nodeCmd.Subcommands, nodeGetCmd)
}

var nodeGetCmd = &cli.Command{
	Name:      "get",
	Usage:     "Show a node detail",
	ArgsUsage: "<user id>",
	Action: ex(func(c *exContext) error {
		type point struct {
			NodeID int             `json:"node_id"`
			Params json.RawMessage `json:"params"`
			Type   string          `json:"type"`
		}
		node, err := c.cnt.NodeShow(c.ctx.Args().First())
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\t%d\n", node.ID)
		fmt.Fprintf(w, "Name\t%s\n", node.Name)

		var obj []interface{}
		json.Unmarshal(node.Points, &obj)
		f := colorjson.NewFormatter()
		f.Indent = 4
		s, err := f.Marshal(obj)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Points\t%s\n", s)
		return w.Flush()
	}),
}
