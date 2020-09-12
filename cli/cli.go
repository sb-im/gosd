package cli

import (
	"flag"
	"fmt"
	"runtime"

	"sb.im/gosd/auth"
	"sb.im/gosd/config"
	"sb.im/gosd/database"
	"sb.im/gosd/storage"

	"miniflux.app/logger"
	"miniflux.app/version"
)

const (
	flagVersionHelp     = "Show application version"
	flagMigrateHelp     = "Run SQL migrations"
	flagCreateAdminHelp = "Create admin user"
	flagDebugModeHelp   = "Show debug logs"
	flagNoAuthHelp      = "use the noauth auther when using quick setup"
)

func Parse() {
	var (
		err             error
		flagVersion     bool
		flagMigrate     bool
		flagCreateAdmin bool
		flagDebugMode   bool
		flagNoAuth      bool
	)

	flag.BoolVar(&flagVersion, "version", false, flagVersionHelp)
	flag.BoolVar(&flagVersion, "v", false, flagVersionHelp)
	flag.BoolVar(&flagMigrate, "migrate", false, flagMigrateHelp)
	flag.BoolVar(&flagCreateAdmin, "create-admin", false, flagCreateAdminHelp)
	flag.BoolVar(&flagDebugMode, "debug", false, flagDebugModeHelp)
	flag.BoolVar(&flagNoAuth, "noauth", false, flagNoAuthHelp)

	flag.Parse()

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

	if flagCreateAdmin {
		createAdmin(store)
		return
	}

	startDaemon(store, opts)
}
