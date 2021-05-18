package config

import (
	"fmt"
	"os"
)

type Config struct {
	plaidCountry	string
	plaidRedirect	string
	plaidClient		string
	plaidSecret		string
	dbUser			string
	dbPassword		string
	dbHost			string
	dbName			string
	port			string
}

// Gets the environment variable configuration necessary to run application
func Get() *Config {
	config := &Config {
		plaidCountry: 	os.Getenv("PLAID_COUNTRY_CODES"),
		plaidRedirect: 	os.Getenv("PLAID_REDIRECT_URI"),
		plaidClient:	os.Getenv("PLAID_CLIENT_ID"),
		plaidSecret: 	os.Getenv("PLAID_SECRET"),
		dbUser:			os.Getenv("DB_USERNAME"),
		dbPassword: 	os.Getenv("DB_PASSWORD"),
		dbHost:			os.Getenv("DB_ACCESSPT"),
		dbName:			os.Getenv("DB_DATABASE"),
		port:			os.Getenv("HOST"),
	}

	return config
}

// Gets the required string to open database instance
func (config *Config) GetDBConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s", 
		config.dbUser, 
		config.dbPassword, 
		config.dbHost, 
		config.dbName,
	)
}

// Gets the port of the server
func (config *Config) GetPort() string {
	return config.port
}