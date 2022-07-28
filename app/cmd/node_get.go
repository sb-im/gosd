package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/TylerBrock/colorjson"
	"github.com/spf13/cobra"
)

func init() {
	nodeCmd.AddCommand(nodeGetCmd)
}

var nodeGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Show a node detail",
	Args:  cobra.ExactArgs(1),
	RunE: ex(func(c *exContext) error {
		type point struct {
			NodeID int             `json:"node_id"`
			Params json.RawMessage `json:"params"`
			Type   string          `json:"type"`
		}
		node, err := c.cnt.NodeShow(c.arg[0])
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\t%s\n", node.UUID)
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
