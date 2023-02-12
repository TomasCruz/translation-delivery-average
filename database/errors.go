package database

import (
	"errors"
)

var (
	errNotPostgresTxErr error = errors.New("not a postgresTx")
)
