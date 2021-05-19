package app

import (
	"database/sql"

	"github.com/elopez00/scale-backend/pkg/config"
	"github.com/elopez00/scale-backend/pkg/database"
)

type App struct {
	DB 		*database.DB
	Config 	*config.Config
}

// app.Get will initialize environment variables and database connection.
// If Get encounters any errors, it will return it.
func Get() (*App, error) {
	configuration := config.Get()
	if db, err := database.Get(configuration.GetDBConnectionString()); err != nil {
		return nil, err
	} else {
		return &App { DB: db, Config: configuration }, nil
	}
}

// app.GetTest will initialize app object given database and configuration
func GetTest(db *sql.DB) *App {
	return &App { 
		DB: database.GetTest(db),
		Config: config.GetTest(),
	}
}