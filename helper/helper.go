package helper

import (
	"net/http/httptest"
	"testing"
	"net/http"
	"time"
	"log"
	"fmt"
	"os"

	"github.com/elopez00/scale-backend/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
		Key: 			os.Getenv("KEY"),
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
		Key: 			os.Getenv("KEY"),
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

func GenerateJWT(keyStr string, id string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: 	id,
		ExpiresAt: 	time.Now().Add(time.Minute * 30).Unix(),
	})

	if token, err := claims.SignedString([]byte(keyStr)); err != nil {
		fmt.Println("Error processing key: ", err.Error())
		return "", err
	} else {
		return token, nil
	}
}

func GetCorsHandler(router *mux.Router) http.Handler {
	c := cors.New(cors.Options {
		AllowedOrigins: []string {
			"http://localhost:5000",
			"http://scale-backend-dev.us-east-1.elasticbeanstalk.com/",
		},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	return handler
}