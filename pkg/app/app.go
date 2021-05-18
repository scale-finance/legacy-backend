package app

import (
	"github.com/elopez00/scale-backend/pkg/config"
	"github.com/elopez00/scale-backend/pkg/database"
)

type App struct {
	DB 		*database.DB
	Config 	*config.Config
}

// gets application
func Get() (*App, error) {
	configuration := config.Get()
	if db, err := database.Get(configuration.GetDBConnectionString()); err != nil {
		return nil, err
	} else {
		return &App { DB: db, Config: configuration, }, nil
	}
	
}