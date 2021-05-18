package main

import (
	"log"

	"github.com/elopez00/scale-backend/cmd/api/router"
	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	} else {
		if app, err := application.Get(); err != nil {
			log.Fatal(err)
		} else {
			defer app.DB.Close()
			srv := server.
				Get().
				WithAddr(app.Config.GetPort()).
				WithHandler(router.Get(app))
	
			log.Println("Starting server at port", app.Config.GetPort())
			if err := srv.Start(); err != nil {
				log.Fatal(err)
			}
		}
	}
}