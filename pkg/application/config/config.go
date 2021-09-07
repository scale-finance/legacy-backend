package config

import (
	"fmt"
)

// TODO make implement dev key for plaid api

type Config struct {
	plaid	 map[string]string
	database map[string]string
	server 	 map[string]string
}

// Get the environment variable configuration necessary to run application
func Get(environment map[string]string) *Config {
	config := &Config {
		plaid: map[string]string {
			"countryCode": 	environment["PLAID_COUNTRY_CODES"],
			"redirectUri": 	environment["PLAID_REDIRECT_URI"],
			"secret": 		environment["PLAID_SECRET"],
			"client": 		environment["PLAID_CLIENT_ID"],
		},
		database: map[string]string {
			"user":		environment["DB_USERNAME"],
			"password": environment["DB_PASSWORD"],
			"host": 	environment["DB_ACCESSPT"],
			"name": 	environment["DB_DATABASE"],
		},
		server: map[string]string {
			"port": environment["HOST"],
			"key":  environment["KEY"],
		},
	}

	return config
}

// GetDBConnectionString retrieves the required string to open database instance
func (config *Config) GetDBConnectionString() string {
	// check if config is using a testing host for the database, if so return a
	// test connection string
	if config.database["host"] == "test" {
		return "test"
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s", 
		config.database["user"],
		config.database["password"],
		config.database["host"],
		config.database["name"],
	)
}

// GetPlaid will retrieve all plaid client details
func (config *Config) GetPlaid() map[string]string {
	return config.plaid
}

// GetServer gets the server details
func (config *Config) GetServer() map[string]string {
	return config.server
}