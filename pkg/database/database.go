package database 

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Client	*sql.DB
}

// Gets the database object after connection, else returns error
func Get(connectionString string) (*DB, error) {
	if db, err := sql.Open("mysql", connectionString); err != nil {
		return nil, err
	} else {
		if err := db.Ping(); err != nil {
			return nil, err
		} 

		return &DB { Client: db, }, nil
	}
}

func (db *DB) Close() error {
	return db.Client.Close()
}