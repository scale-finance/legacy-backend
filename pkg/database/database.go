package database 

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Client	*sql.DB
}

// Recieves a string that will describe its credentials and will return a database object
// If Get encounters any errors, it will return it.
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

func GetTest(client *sql.DB) (*DB) {
	return &DB { Client: client }
}

func (db *DB) Close() error {
	return db.Client.Close()
}