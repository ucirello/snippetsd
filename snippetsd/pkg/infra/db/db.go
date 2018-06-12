package db // import "cirello.io/snippetsd/pkg/infra/db"

import (
	"cirello.io/snippetsd/pkg/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

// Config defines the environment for the database
type Config struct {
	Filename string
}

// Connect dials to the database.
func Connect(config Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", config.Filename)
	if err != nil {
		return nil, errors.E(err)
	}
	return db, nil
}
