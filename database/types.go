package database

import (
	"database/sql"
)

type (
	postgresDB struct {
		db *sql.DB
	}

	postgresTx struct {
		*sql.Tx
	}
)
