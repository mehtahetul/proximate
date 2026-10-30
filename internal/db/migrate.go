package db

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // registers the "pgx5://" database driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/mehtahetul/proximate/internal/db/migrations"
)

// RunMigrations applies every migration file newer than the version
// recorded in the schema_migrations table. On an empty database it runs
// all of them; on an up-to-date database it does nothing.
func RunMigrations() {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("Failed to load migration files: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, databaseURL("pgx5"))
	if err != nil {
		log.Fatalf("Failed to initialise migrations: %v", err)
	}
	defer m.Close()

	// ErrNoChange just means "already up to date" — not a failure.
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	version, dirty, _ := m.Version()
	log.Printf("Database schema at version %d (dirty=%v)", version, dirty)
}
