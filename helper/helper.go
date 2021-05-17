package helper

import (
	"net/http/httptest"
	"testing"
	"net/http"
	"log"
	"os"

	"github.com/elopez00/scale-backend/models"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
)

// initializes environment variables for production
func InitEnv() models.Env {
	// check environemnt
	if os.Getenv("LIVE") != "true" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal(err.Error())
		}
	}

	// get environment variables
	env := models.Env {
		DbUser: 		os.Getenv("DB_USERNAME"),
		DbPass: 		os.Getenv("DB_PASSWORD"),
		DbData: 		os.Getenv("DB_DATABASE"),
		DbAccess: 		os.Getenv("DB_ACCESSPT"),
		PlaidCtry: 		os.Getenv("PLAID_COUNTRY_CODES"),
		PlaidSecret: 	os.Getenv("PLAID_SECRET"),
		PlaidId: 		os.Getenv("PLAID_CLIENT_ID"),
		PlaidRedir: 	os.Getenv("PLAID_REDIRECT_URI"),
		Live: 			os.Getenv("LIVE"),
	}

	return env
}

// initializes environment variables for testing
func InitTestEnv() models.Env {
	if os.Getenv("LIVE") != "true" {
		if err := godotenv.Load(os.ExpandEnv("../.env")); err != nil {
			log.Fatal(err.Error())
		}
	}

	// get environment variables
	env := models.Env {
		DbUser: 		os.Getenv("DB_USERNAME"),
		DbPass: 		os.Getenv("DB_PASSWORD"),
		DbData: 		os.Getenv("DB_DATABASE"),
		DbAccess: 		os.Getenv("DB_ACCESSPT"),
		PlaidCtry: 		os.Getenv("PLAID_COUNTRY_CODES"),
		PlaidSecret: 	os.Getenv("PLAID_SECRET"),
		PlaidId: 		os.Getenv("PLAID_CLIENT_ID"),
		PlaidRedir: 	os.Getenv("PLAID_REDIRECT_URI"),
		Live: 			os.Getenv("LIVE"),
	}

	return env
}

func ExecuteRequest(req *http.Request, router *mux.Router) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func CheckResponse(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}