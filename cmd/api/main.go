package main

import (
	"log"
	"os"

	"github.com/elopez00/scale-backend/cmd/api/router"
	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	// gets environment variables
	if os.Getenv("LIVE") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		} 
	}

	// gets application
	app, err := application.Get()
	if err != nil {
		log.Fatal(err)
	} 

	// defers the close of database
	defer app.DB.Close()

	// creates the server 
	srv := server.
		Get().
		WithAddr(app.Config.GetPort()).
		WithHandler(router.Get(app))
		
	// starts the server on given port
	log.Println("Starting server at port", app.Config.GetPort())
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}