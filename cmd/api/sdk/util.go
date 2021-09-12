package sdk

import (
	"fmt"
	"github.com/elopez00/scale-backend/cmd/api/models"
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

// GetIDFromContext will get the user ID from the applications context
func GetIDFromContext(request *http.Request) string {
	return fmt.Sprintf("%v", request.Context().Value(models.Key("user")))
}