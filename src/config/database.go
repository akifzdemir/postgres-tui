package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
	mu   sync.Mutex
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

func RemoveDb() error {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
		db = nil
		once = sync.Once{}
	}
	return nil
}

func GetDb() (*sql.DB, error) {
	if db == nil {
		return nil, nil
	}
	return db, nil
}
