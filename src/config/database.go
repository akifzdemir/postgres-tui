package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

func ConnectDb(conStr string) (*sql.DB, error) {
	var err error
	once.Do(func() {
		db, err = sql.Open("postgres", conStr)
		if err != nil {
			db = nil
		}

		if err = db.Ping(); err != nil {
			db = nil
		}
	})

	return db, err
}

func GetDb() (*sql.DB, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}
	return db, nil
}
