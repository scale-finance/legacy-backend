package models

import (
	"net/http"
	"encoding/json"
	"log"
)

// for use of json responses
type Response struct {
	Status		int			`json:"status"`
	Message		string		`json:"message"`
	Result		interface{}	`json:"result,omitempty"`
}

// Function that creates and writes a JSON response that was successfully executed. This 
// function will alwayrs return the http status 200 (OK), and has the option to return a
// result which can be any datatype
func CreateResponse(w http.ResponseWriter, message string, result interface{}) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	
	res := Response {
		Status: http.StatusOK,
		Message: message,
		Result: result,
	}

	encoder.Encode(res)
}

// This function is used to create an error JSON response with custom http statuses. 
// In addition to the status, this function also takes in an error that will be logged
// to the system. This error can be nil
func CreateError(w http.ResponseWriter, status int, message string, system error) {
	if system != nil {
		log.Print(system.Error())
	}
	encoder := json.NewEncoder(w)
	
	res := Response {
		Status: status,
		Message: message,
	}
	
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
	encoder.Encode(res)
}

// for use of context keys 
type Key string