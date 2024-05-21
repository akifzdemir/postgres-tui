package config

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDb(conStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
