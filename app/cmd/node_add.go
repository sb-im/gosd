package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	nodeCmd.AddCommand(nodeAddCmd)

	nodeAddCmd.Flags().UintP("team", "t", 0, "Team Id")
	nodeAddCmd.Flags().StringP("username", "u", "", "new username")
}

var nodeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new node",
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
		return c.cnt.NodeCreate(&map[string]interface{}{
			"team_id": mustGetUint(c.ctx.Flags(), "team"),
			"name":    mustGetString(c.ctx.Flags(), "name"),
			"points":  json.RawMessage(points),
		})
	}),
}
