package application

import (
	"github.com/elopez00/scale-backend/pkg/application/config"
	"github.com/elopez00/scale-backend/pkg/application/database"
	"github.com/elopez00/scale-backend/pkg/application/plaid"
)

type App struct {
	// DB is where the database client lives, which is in charge of all database
	// functionalities
	DB 		*database.DB

	// Config is where the application's configuration lies
	Config 	*config.Config

	// Plaid is where the plaid client lives, which is in charge of all bank
	// information retrieval and functionalities
	Plaid	*plaid.Plaid
}

// Get will initialize environment variables and database connection.
// If the function encounters any errors, it will return it.
func Get(environment map[string]string) (*App, error) {
	// get the config
	Config := config.Get(environment)

	// get the database client
	DB, err := database.Get(*Config)
	if err != nil {
		return nil, err
	}

	// get plaid client
	Plaid, err := plaid.Get(*Config)
	if err != nil {
		return nil, err
	}

	return &App { DB: DB, Config: Config, Plaid: Plaid }, nil
}