package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"sb.im/gosd/app/client"
	"sb.im/gosd/app/config"

	"github.com/urfave/cli/v2"
)

func init() {
	userCmd.Subcommands = append(userCmd.Subcommands, userLsCmd)
}

var userLsCmd = &cli.Command{
	Name:  "ls",
	Usage: "ls all users",
	Action: func(c *cli.Context) error {
		cc := client.NewClient(config.Opts().BaseURL)
		users, err := cc.UserIndex()
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tUsername\tTeam\tTeams\tLanguage\tTimezone")

		for _, u := range users {
			fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\t%s\t\n",
				u.ID,
				u.Username,
				u.TeamID,
				func() string {
					arr := make([]string, len(u.Teams))
					for i, t := range u.Teams {
						arr[i] = fmt.Sprintf("%d", t.ID)
					}
					return strings.Join(arr, ",")
				}(),
				u.Language,
				u.Timezone,
			)
		}

		return w.Flush()
	},
}
