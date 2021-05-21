package config

import (
	"fmt"
	"os"
)

// TODO make implement dev key for plaid api

type Config struct {
	plaid			map[string]string
	dbUser			string
	dbPassword		string
	dbHost			string
	dbName			string
	port			string
	key 			string
}

// Gets the environment variable configuration necessary to run application
func Get() *Config {
	config := &Config {
		plaid: map[string]string {
			"country": 		os.Getenv("PLAID_COUNTRY_CODES"),
			"redirectUrl": 	os.Getenv("PLAID_REDIRECT_URI"),
			"secret": 		os.Getenv("PLAID_SECRET"),
			"client": 		os.Getenv("PLAID_CLIENT_ID"),
		},
		dbUser:			os.Getenv("DB_USERNAME"),
		dbPassword: 	os.Getenv("DB_PASSWORD"),
		dbHost:			os.Getenv("DB_ACCESSPT"),
		dbName:			os.Getenv("DB_DATABASE"),
		port:			os.Getenv("HOST"),
		key:			os.Getenv("KEY"),
	}

	return config
}

func GetTest() *Config {
	config := &Config {
		plaid: map[string]string {
			"secret": 		os.Getenv("PLAID_SECRET"),
			"client": 		os.Getenv("PLAID_CLIENT_ID"),
		},
		key: "testkey",
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

// Gets encryption secret key
func (config *Config) GetKey() string {
	return config.key
}

// Gets plaid credentials
func (config *Config) GetPlaid() map[string]string {
	return config.plaid
}