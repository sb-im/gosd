package cmd

import (
	"fmt"
	"runtime"

	"sb.im/gosd/auth"
	"sb.im/gosd/config"
	"sb.im/gosd/cli"
	"sb.im/gosd/database"
	"sb.im/gosd/storage"

	"github.com/spf13/cobra"

	"miniflux.app/version"
	"miniflux.app/logger"
)

var (
	flagVersion     bool
	flagMigrate     bool
	flagDebugMode   bool
	flagNoAuth      bool
)

func init() {
	flags := rootCmd.Flags()
	flags.BoolVarP(&flagVersion, "version", "v", false, "Show application version")
	flags.BoolVar(&flagMigrate, "migrate", false, "Run SQL migrations")
	flags.BoolVar(&flagDebugMode, "debug", false, "Show debug logs")
	flags.BoolVar(&flagNoAuth, "noauth", false, "Use the noauth auther. user.ID == 1")
}

var rootCmd = &cobra.Command{
	Use:   "gosd",
	Short: "gosd cli",
	Long: `
TODO
	`,
	Run: ex(func(cmd *cobra.Command, args []string, d exData) {
		if flagVersion {
			fmt.Printf("gosd %s %s %s %s\n", version.Version, runtime.GOOS, runtime.GOARCH, version.BuildDate)
			return
		}

		if flagNoAuth {
			fmt.Println("=== Enable noauth ===")
			auth.SetAuthMethod(auth.NoAuth)
		}

		parse := config.NewParser()
		opts, err := parse.ParseEnvironmentVariables()
		if err != nil {
			panic(err)
		}

		if flagDebugMode || opts.HasDebugMode() {
			logger.EnableDebug()
		}

		db, err := database.NewConnectionPool(
			opts.DatabaseURL(),
			opts.DatabaseMinConns(),
			opts.DatabaseMaxConns(),
		)

		if err != nil {
			panic(err)
		}

		if flagMigrate {
			database.Migrate(db)
			return
		}

		store := storage.NewStorage(db)
		cli.StartDaemon(store, opts)
	}, exConfig{}),
}
