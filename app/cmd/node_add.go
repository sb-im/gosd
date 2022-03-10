package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"
)

func init() {
	nodeCmd.Subcommands = append(nodeCmd.Subcommands, nodeAddCmd)
}

var nodeAddCmd = &cli.Command{
	Name:  "add",
	Usage: "Create a new node",
	Flags: []cli.Flag{
		&cli.UintFlag{Name: "team"},
		&cli.StringFlag{Name: "name"},
	},
	ArgsUsage: "<point path>",
	Action: ex(func(c *exContext) error {
		path := c.ctx.Args().First()
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		return c.cnt.NodeCreate(&map[string]interface{}{
			"team_id": c.ctx.Uint("team"),
			"name":    c.ctx.String("name"),
			"points":  json.RawMessage(data),
		})
	}),
}
