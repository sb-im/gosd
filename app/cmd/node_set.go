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
		&cli.StringFlag{Name: "points"},
	},
	ArgsUsage: "<node id>",
	Action: ex(func(c *exContext) error {
		path := c.ctx.String("points")
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		return c.cnt.NodeUpdate(c.ctx.Args().First(), &map[string]interface{}{
			"team_id": c.ctx.Uint("team"),
			"name":    c.ctx.String("name"),
			"points":  json.RawMessage(data),
		})
	}),
}
