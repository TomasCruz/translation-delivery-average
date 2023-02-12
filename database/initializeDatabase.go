package database

import (
	"database/sql"
	"log"

	"github.com/TomasCruz/translation-delivery-average/service"
	"github.com/pkg/errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// file driver for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// InitializeDatabase does DB migrations and verifies DB accessibility
func InitializeDatabase(dbString string) (pDb service.Database, err error) {
	var db *sql.DB
	if db, err = sql.Open("postgres", dbString); err != nil {
		err = errors.WithStack(err)
		return
	}

	var driver database.Driver
	if driver, err = postgres.WithInstance(db, &postgres.Config{}); err != nil {
		err = errors.WithStack(err)
		return
	}

	var m *migrate.Migrate
	if m, err = migrate.NewWithDatabaseInstance("file://database/migrations", "", driver); err != nil {
		err = errors.WithStack(err)
		return
	}

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			err = nil
		} else {
			return
		}
	}

	sErr, dbErr := m.Close()
	if sErr != nil {
		log.Fatal(sErr)
	} else if dbErr != nil {
		log.Fatal(dbErr)
	}

	if db, err = sql.Open("postgres", dbString); err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	pDb = postgresDB{db: db}
	return
}
