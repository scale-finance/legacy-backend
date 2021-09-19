package sdk

import (
	"fmt"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"log"
	"net/http"
	"strings"
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

// GetPlaidErrorCode will get the error code from the error message and return it as a string
func GetPlaidErrorCode(err error) string {
	errorMessage := err.Error()

	// first get the index of the substring code
	start := strings.Index(errorMessage, ", code: ") + 8

	// get the end by creating a substring and getting the index of the first comma
	end := strings.Index(errorMessage[start:], ", ") + start

	// return the substring with the window of indeces
	return errorMessage[start:end]
}