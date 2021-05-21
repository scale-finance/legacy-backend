package app

import (
	"database/sql"
	"log"

	"github.com/elopez00/scale-backend/pkg/config"
	"github.com/elopez00/scale-backend/pkg/database"
	"github.com/elopez00/scale-backend/pkg/plaid"
)

type App struct {
	DB 		*database.DB
	Config 	*config.Config
	Plaid	*plaid.Plaid
}

// app.Get will initialize environment variables and database connection.
// If Get encounters any errors, it will return it.
func Get() (*App, error) {
	configuration := config.Get()
	plaidConfig := configuration.GetPlaid()
	if db, err := database.Get(configuration.GetDBConnectionString()); err != nil {
		return nil, err
	} else {
		plaid, err := plaid.Get(plaidConfig)
		if err != nil {
			return nil, err
		}

		return &App { DB: db, Config: configuration, Plaid: plaid }, nil
	}
}

// app.GetTest will initialize app object given database and configuration
func GetTest(db *sql.DB) *App {
	configuration := config.GetTest()
	plaidOptions := configuration.GetPlaid()
	if p, err := plaid.Get(plaidOptions); err != nil {
		log.Print("Getting app")
		panic(err.Error())
	} else {
		return &App { 
			DB: database.GetTest(db),
			Config: configuration,
			Plaid: p,
		}
	}
}