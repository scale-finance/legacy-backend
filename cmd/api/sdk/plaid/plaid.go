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

// Returns the plaid token from authentication token. If in any case there is an error with
// the link token or the user's connection, it will return a json response error to the frontend
func GetPlaidToken(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")		

		// creates token configuration
		tokenConfig := plaid.LinkTokenConfigs {
			User: &plaid.LinkTokenUser {
				ClientUserID: fmt.Sprintf("%v", r.Context().Value(models.Key("user"))),
			},
			ClientName: 	"Scale",
			Products: 		[]string{"auth", "transactions"},
			CountryCodes: 	[]string{"US"},
			Language:		"en",
			Webhook:		app.Plaid.RedirectURL,
		}

		// calls on the app's plaid client and creates a link token with the configuration
		// provided by the tokenConfig struct. If for whatever reason the client fails, it 
		// will return a json resposne reflecting this issue
		tokenResponse, err := app.Plaid.Client.CreateLinkToken(tokenConfig)
		if err != nil {
			log.Println("Error Loading Client:", err)
			json.NewEncoder(w).Encode(models.Response {
				Type: "Link Token",
				Status: 1,
				Message: "Failure to load client",
			})
		}
		
		// when successful returns a result response
		json.NewEncoder(w).Encode(models.Response {
			Result: tokenResponse.LinkToken,
			Type: "Link Token",
			Message: "Successfully recieved link token from plaid",
			Status: 0,
		})
	}
}

// this function takes care of creating the permanent access token that will be 
// stored in the database for cross-platform connection to users' bank. If for
// whatever reason there is a problem with the client or public token, their 
// are json responses and logs that will adequately reflect all issues
// ! Not tested yet
func CreateAccessToken(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var token models.Tkn
		json.NewDecoder(r.Body).Decode(&token)
		
		// Creates the permanent token using the public token it gets from the frontend's
		// request body
		res, err := app.Plaid.Client.ExchangePublicToken(token.Public)
		if err != nil {
			log.Println("Failed create access token", err)
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Access Token",
				Message: "Link token error",
			}); return
		}

		// creates a user with an id extracted from authentication cookie for 
		// later use in the creation of the row containing the permanent token
		user := models.User {
			Id: fmt.Sprintf("%v", r.Context().Value(models.Key("user"))),
		}

		// handles failures in the addition of tokens to the database and reflects
		// any success or failure in json response/server logs
		if err = user.AddToken(app, res.AccessToken, res.ItemID); err != nil {
			log.Println("Failed to add token to DB:", err)
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Access Token",
				Message: "Failed to create access token",
			}); return
		} else {
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Access Token",
				Message: "Access Token successfuly created",
			})
		}
	}
}