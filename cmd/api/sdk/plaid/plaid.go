package plaid

import (
	"log"
	"fmt"
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
				Type: "Link Token",
				Status: 1,
				Message: "Failure to load client",
			})
		} else {
			// create configuration
			tokenConfig := plaid.LinkTokenConfigs {
				User: &plaid.LinkTokenUser {
					ClientUserID: fmt.Sprintf("%v", r.Context().Value(models.Key("user"))),
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
					Type: "Link Token",
					Status: 1,
					Message: "Failure to load client",
				})
			} else {
				json.NewEncoder(w).Encode(models.Response {
					Result: tokenResponse.LinkToken,
					Type: "Link Token",
					Message: "Successfully recieved link token from plaid",
					Status: 0,
				})
			}
		}
	}
}