package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"strconv"

	users "sb.im/gosd/model"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(usersCmd)
}

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Users management utility",
	Long:  `Users management utility.`,
	Args:  cobra.NoArgs,
}

func printUsers(usrs []*users.User) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tUsername\tLanguage\tTimezone")

	for _, u := range usrs {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n",
			u.ID,
			u.Username,
			u.Language,
			u.Timezone,
		)
	}

	w.Flush()
}

func parseUsernameOrID(arg string) (username string, id int64) {
	id64, err := strconv.ParseInt(arg, 10, 0)
	if err != nil {
		return arg, 0
	}
	return "", id64
}
