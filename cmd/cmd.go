package cmd

// Execute executes the commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
