package database

import (
	"context"
	"database/sql"
)

func (pDB postgresDB) newTransaction() (tx postgresTx, err error) {
	var sqlTx *sql.Tx
	sqlTx, err = pDB.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return
	}

	tx = postgresTx{sqlTx}
	return
}

func (pTx postgresTx) commitOrRollbackOnError(err *error) {
	if *err != nil {
		pTx.Rollback()
	} else {
		pTx.Commit()
	}
}
