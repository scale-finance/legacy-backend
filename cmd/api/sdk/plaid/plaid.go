package Plaid

import (
	"log"
	"net/http"
	"encoding/json"

	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/julienschmidt/httprouter"
	"github.com/plaid/plaid-go/plaid"
)

// Returns the plaid token given a body with id
func GetPlaidToken(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		plaidConfig := app.Config.GetPlaid()
		clientOptions := plaid.ClientOptions {
			ClientID:		plaidConfig["client"],
			Secret:			plaidConfig["secret"],
			Environment:	plaid.Sandbox,
			HTTPClient:		&http.Client{},
		}

		if client, err := plaid.NewClient(clientOptions); err != nil {
			log.Println("Error Loading Client")
			json.NewEncoder(w).Encode(models.Response {
				Type: "Plaid Token",
				Status: 1,
				Message: err.Error(),
			})
		} else {
			// TODO retrieve user id from auth token
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
				Webhook:		plaidConfig["redirectUrl"],
			}

			// create and return token
			if tokenResponse, err := client.CreateLinkToken(tokenConfig); err != nil {
				log.Println("Error Loading Client")
				json.NewEncoder(w).Encode(models.Response {
					Type: "Plaid Token",
					Status: 1,
					Message: err.Error(),
				})
			} else {
				type response = struct { LinkToken string `json:"linktoken"` }
				json.NewEncoder(w).Encode(response {
					LinkToken: tokenResponse.LinkToken,
				})
			}
		}
	}
}