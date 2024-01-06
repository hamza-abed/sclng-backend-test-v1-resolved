package util

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

func RunCreateSQLStatement(sql string, db *sql.DB, logger logrus.FieldLogger) error {
	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
	tx.Commit()
	if err != nil {
		panic(err)
	}
	return err
}
