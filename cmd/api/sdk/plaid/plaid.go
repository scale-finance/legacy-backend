package plaid

import (
	"encoding/json"
	"fmt"
	"time"
	"net/http"

	"github.com/elopez00/scale-backend/cmd/api/models"
	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/julienschmidt/httprouter"
	"github.com/plaid/plaid-go/plaid"
)

// Returns the plaid token from authentication token. If in any case there is an error with
// the link token or the user's connection, it will return a json response error to the frontend
func GetPlaidToken(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		// creates token configuration
		tokenConfig := plaid.LinkTokenConfigs{
			User: &plaid.LinkTokenUser{
				ClientUserID: fmt.Sprintf("%v", r.Context().Value(models.Key("user"))),
			},
			ClientName:   "Scale",
			Products:     []string{"auth", "transactions"},
			CountryCodes: []string{"US"},
			Language:     "en",
			Webhook:      app.Plaid.RedirectURL,
		}

		// calls on the app's plaid client and creates a link token with the configuration
		// provided by the tokenConfig struct. If for whatever reason the client fails, it
		// will return a json resposne reflecting this issue
		tokenResponse, err := app.Plaid.Client.CreateLinkToken(tokenConfig)
		if err != nil {
			msg := "Failure to load client"
			models.CreateError(w, http.StatusBadGateway, msg, err)
			return
		}

		// when successful returns a result response
		msg := "Successfully recieved link token from plaid"
		models.CreateResponse(w, msg, tokenResponse.LinkToken)
	}
}

// this function takes care of creating the permanent access token that will be
// stored in the database for cross-platform connection to users' bank. If for
// whatever reason there is a problem with the client or public token, their
// are json responses and logs that will adequately reflect all issues
// ! Not tested yet
func CreateAccessToken(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var token models.Token
		json.NewDecoder(r.Body).Decode(&token)

		// Creates the permanent token using the public token it gets from the frontend's
		// request body
		res, err := app.Plaid.Client.ExchangePublicToken(token.Value)
		if err != nil {
			msg := "Failure to excahnge link token"
			models.CreateError(w, http.StatusBadGateway, msg, err)
			return
		}

		// sets value of token to response values
		token.Value = res.AccessToken
		token.Id = res.ItemID

		// gets the user id extracted from authentication cookie for later
		// use in the creation of the row containing the permanent token
		userId := fmt.Sprintf("%v", r.Context().Value(models.Key("user")))

		// handles failures in the addition of tokens to the database and reflects
		// any success or failure in json response/server logs
		if err = token.Add(app, userId); err != nil {
			msg := "Failure to create access token"
			models.CreateError(w, http.StatusBadGateway, msg, err)
			return
		} 
		
		msg := "Access token created successfully"
		models.CreateResponse(w, msg, nil)
	}
}

// This function will get transactions from the past 2 months from all bank accounts 
// affiliated with the user. If there is an error with the database retrieval or the
// plaid client call, this will be reflected in the json response accordingly.
func GetTransactions(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// create the user with the id obtained from middleware context
		user := models.User { Id: fmt.Sprintf("%v", r.Context().Value(models.Key("user"))) }

		tokens, err := user.GetTokens(app)
		if err != nil {
			msg := "There was an error retrieving tokens from database affiliated with user"
			models.CreateError(w, http.StatusBadGateway, msg, err)
		}

		const iso8601TimeFormat = "2006-01-02"
		startDate := time.Now().Add(-30 * 24 * time.Hour).Format(iso8601TimeFormat)
		endDate := time.Now().Format(iso8601TimeFormat)

		var transactionHistory []plaid.Transaction
		for i := range tokens {
			res, err := app.Plaid.Client.GetTransactions(tokens[i].Value, startDate, endDate)
			if err != nil {
				msg := "Failed to retrieve tokens from Plaid client"
				models.CreateError(w, http.StatusBadGateway, msg, err)
				return
			}

			transactionHistory = append(transactionHistory, res.Transactions...)
		}

		json.NewEncoder(w).Encode(models.Response {
			Status: http.StatusOK,
			Message: "Successfully retrieved transactions from all bank accounts",
			Result: transactionHistory,
		})
	}
}