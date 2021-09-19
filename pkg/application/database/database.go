package database

import (
	"database/sql"
	"log"

	"github.com/elopez00/scale-backend/pkg/application/config"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	// Client is the database client
	Client	*sql.DB

	// Mock will only be available in a test environment and is used for mocking database functions
	Mock	sqlmock.Sqlmock
}

// Get receives a string that will describe its credentials and will return a database object
// If the function encounters any errors, it will return it.
func Get(config config.Config) (*DB, error) {
	// get the connection string
	connectionString := config.GetDBConnectionString()

	// if the connection string is test, return a test database and mock instance
	if connectionString == "test" {
		db, mock, _ := sqlmock.New()
		return &DB { Client: db, Mock: mock }, nil
	}

	// else, we get the database described by the environment
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Print(connectionString)
		return nil, err
	}

	// if we can't successfully ping the server than we return an error
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB { Client: db, Mock: nil }, nil
}

// Close will close the database client
func (db *DB) Close() error {
	return db.Client.Close()
}