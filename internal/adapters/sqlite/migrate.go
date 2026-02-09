package sqlite

import (
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var schema string

func migrate(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
