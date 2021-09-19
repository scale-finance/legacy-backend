package main

import (
	"log"
	"os"
	"strings"

	"github.com/elopez00/scale-backend/cmd/api/router"
	"github.com/elopez00/scale-backend/pkg/application"
	"github.com/elopez00/scale-backend/pkg/application/database"
	"github.com/elopez00/scale-backend/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	// gets environment variables
	environment, err := godotenv.Read()
	if err != nil {
		// if for whatever reason godotenv fails, try to load from the system environment
		// and create a map with its values
		environment := make(map[string]string)
		for _, item := range os.Environ() {
			log.Println("Environment Line: ", item)
			splits := strings.Split(item, "=")
			key := strings.Trim(splits[0], " ")
			val := strings.Trim(splits[1], " ")
			log.Println("Key: ", key, " Value: ", val)
			environment[key] = val
		}

		log.Println("Environment map: ", environment)
		
		// if there were no elements then either return the error again
		if len(environment) == 0 {
			log.Fatal("Failed go get environment file:", err)
		}
	}

	log.Println(environment)

	// gets application
	app, err := application.Get(environment)
	if err != nil {
		log.Fatal(err)
	}

	// defers the close of database
	defer func(DB *database.DB) {
		err := DB.Close()
		if err != nil {
			log.Println("Failed to close database instance", err)
		}
	} (app.DB)

	// creates the server
	port := app.Config.GetServer()["port"]
	srv := server.
		Get().
		WithAddr(port).
		WithHandler(router.Get(app))

	// starts the server on given port
	log.Println("Starting server at port", port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
