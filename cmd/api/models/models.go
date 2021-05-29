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
	Type 		string 		`json:"type"`
	Result		interface{}	`json:"result,omitempty"`
}

func CreateResponse(w http.ResponseWriter, message string, result interface{}) {
	encoder := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	
	res := Response {
		Status: http.StatusOK,
		Message: message,
		Result: result,
	}

	encoder.Encode(res)
}

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