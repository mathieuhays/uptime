package uptime

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "sql/schema"); err != nil {
		return err
	}

	return nil
}
