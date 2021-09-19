package main

import (
	"github.com/elopez00/scale-backend/cmd/api/router"
	"github.com/elopez00/scale-backend/pkg/application"
	"github.com/elopez00/scale-backend/pkg/application/database"
	"github.com/elopez00/scale-backend/pkg/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// gets environment variables
	environment, err := godotenv.Read()
	if err != nil {
		log.Fatal("Failed go get environment file:", err)
	}

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
