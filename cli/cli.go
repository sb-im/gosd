package cli

import (
	"errors"
	"flag"
	"fmt"
	"runtime"

	"sb.im/gosd/config"
	"sb.im/gosd/database"
	"sb.im/gosd/storage"

	"miniflux.app/version"
)

const (
	flagVersionHelp     = "Show application version"
	flagMigrateHelp     = "Run SQL migrations"
	flagCreateAdminHelp = "Create admin user"
	flagDebugModeHelp   = "Show debug logs"
)

func Parse() (*config.Options, error) {
	var (
		err             error
		flagVersion     bool
		flagMigrate     bool
		flagCreateAdmin bool
		flagDebugMode   bool
	)

	flag.BoolVar(&flagVersion, "version", false, flagVersionHelp)
	flag.BoolVar(&flagVersion, "v", false, flagVersionHelp)
	flag.BoolVar(&flagMigrate, "migrate", false, flagMigrateHelp)
	flag.BoolVar(&flagCreateAdmin, "create-admin", false, flagCreateAdminHelp)
	flag.BoolVar(&flagDebugMode, "debug", false, flagDebugModeHelp)

	flag.Parse()

	if flagVersion {
		fmt.Printf("gosd %s %s %s %s\n", version.Version, runtime.GOOS, runtime.GOARCH, version.BuildDate)
		return nil, errors.New("=== show version ===")
	}

	parse := config.NewParser()
	opts, err := parse.ParseEnvironmentVariables()
	if err != nil {
		return opts, err
	}

	if flagMigrate {
		db, err := database.NewConnectionPool(
			opts.DatabaseURL(),
			opts.DatabaseMinConns(),
			opts.DatabaseMaxConns(),
		)

		if err != nil {
			return opts, err
		}

		database.Migrate(db)
		return opts, errors.New("=== end migrate ===")
	}

	if flagCreateAdmin {
		db, err := database.NewConnectionPool(
			opts.DatabaseURL(),
			opts.DatabaseMinConns(),
			opts.DatabaseMaxConns(),
		)

		if err != nil {
			return opts, err
		}

		store := storage.NewStorage(db)

		createAdmin(store)
		return opts, errors.New("=== end createAdmin ===")
	}

	return opts, err
}
