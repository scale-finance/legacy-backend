package main

import (
	"fmt"
	"net/http"

	"github.com/elopez00/scale-backend/database"
	"github.com/elopez00/scale-backend/helper"
	"github.com/elopez00/scale-backend/sdk"
	"github.com/gorilla/mux"
)

func main() { 
	env := helper.InitEnv() // gets env variables
	router := mux.NewRouter() // creates new mux router
	db := database.Connect(router, env) // initializes database
	
	// handlers
	sdk.Connect(router, env)
	handler := helper.GetCorsHandler(router)

	// close db after
	defer db.Close()
	
	// listen to port 5000
	fmt.Println("Listening on port 5000")
	http.ListenAndServe(":5000", handler)
}
