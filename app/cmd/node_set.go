package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	nodeCmd.AddCommand(nodeSetCmd)
	nodeSetCmd.Flags().Uint("team", 0, "Team ID")
	nodeSetCmd.Flags().String("name", "", "Node Name")
	nodeSetCmd.Flags().String("uuid", "", "Node uuid")
	nodeSetCmd.Flags().String("id", "", "id")
}

var nodeSetCmd = &cobra.Command{
	Use:   "set <point path> <point path> ...",
	Short: "Update a node",
	RunE: ex(func(c *exContext) error {
		ps := make([]json.RawMessage, len(c.arg))
		if len(c.arg) > 0 {
			for i, v := range c.arg {
				f, err := os.Open(v)
				if err != nil {
					panic(err)
				}
				data, err := ioutil.ReadAll(f)
				if err != nil {
					panic(err)
				}

				ps[i] = json.RawMessage(data)
			}
		}
		points, err := json.Marshal(ps)
		if err != nil {
			panic(err)
		}

		return c.cnt.NodeUpdate(mustGetString(c.ctx.Flags(), "id"), &map[string]interface{}{
			"team_id": mustGetUint(c.ctx.Flags(), "team"),
			"name":    mustGetString(c.ctx.Flags(), "name"),
			"uuid":    mustGetString(c.ctx.Flags(), "uuid"),
			"points":  json.RawMessage(points),
		})
	}),
}
