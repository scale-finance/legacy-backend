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

	// initializes database
	db := database.Connect(router, env) 
	
	// handlers
	router.HandleFunc("/", sdk.Greet)
	router.HandleFunc("/getPlaidToken", sdk.GetLinkToken).Methods("GET")
	
	// close db after
	defer db.Close()
	
	// listen to port 5000
	fmt.Println("Listening on port 5000")
	http.ListenAndServe(":5000", router)
}
