package database

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL, migrationPath string, logger *slog.Logger) error {
	if databaseURL == "" {
		return fmt.Errorf("Database URL is not set")
	}

	sourceUrl := fmt.Sprintf("file://%s", migrationPath)

	m, err := migrate.New(sourceUrl, databaseURL)
	if err != nil {
		return fmt.Errorf("initialising migrator: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("reading migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf(
			"database is in a dirty state at version %d — fix manually with: migrate force %d",
			version, version,
		)
	}

	logger.Info("Running database migrations", slog.Uint64("current_version", uint64(version)), slog.String("migrations_path", migrationPath))

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("Database schema is up to date")
			return nil
		}

		return fmt.Errorf("Applying migrations: %w", err)
	}

	newVersion, _, _ := m.Version()
	logger.Info("Migration applied successfully", slog.Uint64("version", uint64(newVersion)))

	return nil
}
