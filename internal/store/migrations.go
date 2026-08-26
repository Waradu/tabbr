package store

import (
	"database/sql"
	"fmt"
	"time"
)

type migration struct {
	version int
	apply   func(*sql.Tx) error
}

var migrations = []migration{
	{
		version: 1,
		apply: func(tx *sql.Tx) error {
			_, err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS commands (
					command TEXT PRIMARY KEY,
					use_count INTEGER NOT NULL DEFAULT 1,
					last_used_at INTEGER NOT NULL
				);

				CREATE TABLE IF NOT EXISTS exclusions (
					pattern TEXT PRIMARY KEY,
					created_at INTEGER NOT NULL
				);
			`)
			return err
		},
	},
	{
		version: 2,
		apply: func(tx *sql.Tx) error {
			return seedDefaultExclusions(tx, []string{
				"cd *",
				"ls *",
				"*access_key*",
				"*api_key*",
				"*secret_key*",
				"*private_key*",
				"*bearer *",
				"*--api-key*",
				"*--token*",
				"*--password*",
				"*--secret*",
			})
		},
	},
	{
		version: 3,
		apply: func(tx *sql.Tx) error {
			_, err := tx.Exec(`
					CREATE INDEX IF NOT EXISTS commands_last_used_at_idx
					ON commands(last_used_at)
				`)
			return err
		},
	},
}

func migrate(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at INTEGER NOT NULL
		)
	`); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	previousVersion := 0
	for _, migration := range migrations {
		if migration.version <= previousVersion {
			return fmt.Errorf("migrations are not ordered at version %d", migration.version)
		}
		previousVersion = migration.version

		if err := applyMigration(db, migration); err != nil {
			return fmt.Errorf("apply migration %d: %w", migration.version, err)
		}
	}

	return nil
}

func applyMigration(db *sql.DB, migration migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.Exec(`
		INSERT INTO schema_migrations (version, applied_at)
		VALUES (?, ?)
		ON CONFLICT(version) DO NOTHING
	`, migration.version, time.Now().Unix())
	if err != nil {
		return err
	}

	inserted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if inserted == 0 {
		return tx.Commit()
	}

	if err := migration.apply(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func seedDefaultExclusions(tx *sql.Tx, patterns []string) error {
	for _, pattern := range patterns {
		if _, err := tx.Exec(`
			INSERT INTO exclusions (pattern, created_at)
			VALUES (?, ?)
			ON CONFLICT(pattern) DO NOTHING
		`, pattern, time.Now().Unix()); err != nil {
			return err
		}

		if _, err := tx.Exec(`
			DELETE FROM commands
			WHERE lower(command) GLOB lower(?)
		`, pattern); err != nil {
			return err
		}
	}

	return nil
}
