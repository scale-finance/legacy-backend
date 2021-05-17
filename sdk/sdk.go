package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/elopez00/scale-backend/middleware"
	"github.com/elopez00/scale-backend/models"
	"github.com/plaid/plaid-go/plaid"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
)

var env models.Env

func Connect(router *mux.Router, environment models.Env) {
	env = environment

	router.HandleFunc("/v0/getPlaidToken", GetLinkToken).Methods("GET")
	router.HandleFunc("/", Greet)
}

// greets user at home directory
func Greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

// Returns the plaid token given a body with id
func GetLinkToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if authenticated := middleware.IsAuthorized(r, env); !authenticated {
		json.NewEncoder(w).Encode(models.Response {
			Status: 1,
			Type: "Authentication",
			Message: "Client not authenticated",
		})
	}

	live := os.Getenv("LIVE")

	// check to see if application is live
	if live != "true" {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Fprint(w, "[ERROR]: ", err.Error())
			panic(err.Error())
		}
	}

	// options given for the plaid client
	clientOptions := plaid.ClientOptions {
		ClientID:		os.Getenv("PLAID_CLIENT_ID"),
		Secret:			os.Getenv(("PLAID_SECRET")),
		Environment:	plaid.Sandbox,
		HTTPClient: 	&http.Client{},
	}

	// create a new client
	if client, err := plaid.NewClient(clientOptions); err != nil {
		fmt.Fprint(w, "[ERROR]: ", err.Error())
		panic(err.Error())
	} else {
		// extract user from body
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
		// create configuration
		tokenConfig := plaid.LinkTokenConfigs {
			User: &plaid.LinkTokenUser {
				ClientUserID: user.Id,
			},
			ClientName: 	"Scale",
			Products: 		[]string{"auth", "transactions"},
			CountryCodes: 	[]string{"US"},
			Language:		"en",
			Webhook:		os.Getenv("PLAID_REDIRECT_URI"),
		}

		// create and return token
		if tokenResponse, err := client.CreateLinkToken(tokenConfig); err != nil {
			fmt.Fprint(w, "[ERROR]: ", err.Error())
			panic(err.Error())
		} else {
			type response = struct { LinkToken string `json:"linktoken"` }
			json.NewEncoder(w).Encode(response {
				LinkToken: tokenResponse.LinkToken,
			})
		}
	}
}