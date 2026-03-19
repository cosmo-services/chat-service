package pkg

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func RunPostgresMigrations(db *sql.DB, migrationPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migration driver error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migration init error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
