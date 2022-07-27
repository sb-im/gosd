package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gosd",
	Short: "SuperDock Cloud Service",
	Long: `StrawBerry Innovation
						SuperDock Cloud Service
						https://sb.im`,
	PreRun: func(c *cobra.Command, args []string) {
		if verbose, _ := c.Flags().GetBool("verbose"); verbose {
			log.SetReportCaller(true)
			log.SetLevel(log.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		Daemon()
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
