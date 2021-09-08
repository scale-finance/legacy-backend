package sdk

import (
	"log"
	"net/http"
)

// CloseBody utility function that repetitive closing handling
func CloseBody(request *http.Request) {
	err := request.Body.Close()
	if err != nil {
		log.Println("Failed to close the body")
	}
}