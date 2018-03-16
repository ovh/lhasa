package db

import (
	"database/sql"

	"github.com/gobuffalo/packr"
	"github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

var migrations = &migrate.PackrMigrationSource{
	Box: packr.NewBox("../../migrations"),
}

func init() {
	migrate.SetTable("ovh_sql_schema_migrations")
}

// MigrateUp run sql-migrate migrations
func MigrateUp(db *sql.DB, log *logrus.Logger) error {
	count, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	log.Warnf("Applied %d migrations", count)
	return nil
}

// MigrateDown run sql-migrate migrations
func MigrateDown(db *sql.DB, log *logrus.Logger) error {
	count, err := migrate.Exec(db, "postgres", migrations, migrate.Down)
	if err != nil {
		return err
	}
	log.Warnf("Removed %d migrations", count)
	return nil
}
