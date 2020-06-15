package database

import (
	"database/sql"
	"fmt"
	"strconv"

	"miniflux.app/logger"
)

const schemaVersion = 1

// Migrate executes database migrations.
func Migrate(db *sql.DB) {
	var currentVersion int
	db.QueryRow(`SELECT version FROM schema_version`).Scan(&currentVersion)

	fmt.Println("Current schema version:", currentVersion)
	fmt.Println("Latest schema version:", schemaVersion)

	for version := currentVersion + 1; version <= schemaVersion; version++ {
		fmt.Println("Migrating to version:", version)

		tx, err := db.Begin()
		if err != nil {
			logger.Fatal("[Migrate] %v", err)
		}

		rawSQL := SqlMap["schema_version_"+strconv.Itoa(version)]
		// fmt.Println(rawSQL)
		_, err = tx.Exec(rawSQL)
		if err != nil {
			tx.Rollback()
			logger.Fatal("[Migrate] %v", err)
		}

		if _, err := tx.Exec(`delete from schema_version`); err != nil {
			tx.Rollback()
			logger.Fatal("[Migrate] %v", err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_version (version) VALUES ($1)`, version); err != nil {
			tx.Rollback()
			logger.Fatal("[Migrate] %v", err)
		}

		if err := tx.Commit(); err != nil {
			logger.Fatal("[Migrate] %v", err)
		}
	}
}

// IsSchemaUpToDate checks if the database schema is up to date.
func IsSchemaUpToDate(db *sql.DB) error {
	var currentVersion int
	db.QueryRow(`SELECT version FROM schema_version`).Scan(&currentVersion)

	if currentVersion != schemaVersion {
		return fmt.Errorf(`the database schema is not up to date: current=v%d expected=v%d`, currentVersion, schemaVersion)
	}
	return nil
}
