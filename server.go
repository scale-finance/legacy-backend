package main

import (
	"fmt"
	"net/http"

	"github.com/elopez00/scale-backend/database"
	"github.com/elopez00/scale-backend/sdk"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	db := database.Connect(router)
	
	router.HandleFunc("/", sdk.Greet)
	router.HandleFunc("/getPlaidToken", sdk.GetLinkToken).Methods("GET")
	
	defer db.Close()
	
	fmt.Println("Listening on port 5000")
	http.ListenAndServe(":5000", router)
}
