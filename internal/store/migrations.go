package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func migrate(db *sql.DB) error {
	ctx := context.Background()
	files, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return err
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, files)
	if err != nil {
		return err
	}
	version, err := provider.GetDBVersion(ctx)
	if err != nil {
		return err
	}
	sources := provider.ListSources()
	latest := sources[len(sources)-1].Version
	if version > latest {
		return fmt.Errorf("database schema version %d is newer than supported version %d; upgrade Tabbr", version, latest)
	}
	_, err = provider.Up(ctx)
	return err
}
