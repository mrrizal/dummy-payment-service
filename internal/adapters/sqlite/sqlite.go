package sqlite

import (
	"database/sql"

	"github.com/XSAM/otelsql"
	_ "github.com/mattn/go-sqlite3"
)

func New(dsn string) (*sql.DB, error) {
	db, err := otelsql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// run migration
	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
