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

		// create configuration
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

		// create and return token
		if tokenResponse, err := app.Plaid.Client.CreateLinkToken(tokenConfig); err != nil {
			log.Println("Error Loading Client", err)
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

func CreateAccessToken(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var token models.Tkn
		json.NewDecoder(r.Body).Decode(&token)
		
		res, err := app.Plaid.Client.ExchangePublicToken(token.Public)
		if err != nil {
			log.Println("Failed create access token", err)
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Access Token",
				Message: "Link token error",
			}); return
		}

		user := models.User {
			Id: fmt.Sprintf("%v", r.Context().Value(models.Key("user"))),
		}

		if err = user.AddToken(app, res.AccessToken, res.ItemID); err != nil {
			log.Println("Failed to add token to DB")
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