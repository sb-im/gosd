package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"
)

func init() {
	nodeCmd.Subcommands = append(nodeCmd.Subcommands, nodeSetCmd)
}

var nodeSetCmd = &cli.Command{
	Name:  "set",
	Usage: "Update a node",
	Flags: []cli.Flag{
		&cli.UintFlag{Name: "team"},
		&cli.StringFlag{Name: "name"},
		&cli.StringFlag{Name: "id"},
	},
	ArgsUsage: "<point path> <point path> ...",
	Action: ex(func(c *exContext) error {
		ps := make([]json.RawMessage, c.ctx.Args().Len())
		if c.ctx.Args().Len() > 0 {
			for i, v := range c.ctx.Args().Slice() {
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

		return c.cnt.NodeUpdate(c.ctx.String("id"), &map[string]interface{}{
			"team_id": c.ctx.Uint("team"),
			"name":    c.ctx.String("name"),
			"points":  json.RawMessage(points),
		})
	}),
}
