package plaid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
func ExchangePublicToken(app *application.App) httprouter.Handle {
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
// TODO make function get transactions that will be seen with budget restrictions
func GetTransactions(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// create the user with the id obtained from middleware context
		userId := fmt.Sprintf("%v", r.Context().Value(models.Key("user"))) 

		// gets all tokens affiliated with user
		tokens, err := models.GetTokens(app, userId)
		if err != nil {
			msg := "There was an error retrieving tokens from database affiliated with user"
			models.CreateError(w, http.StatusBadGateway, msg, err)
		}

		// determines the ending and starting dates to retrieve transactions
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

		msg := "Successfully retrieved transactions from all bank accounts"
		models.CreateResponse(w, msg, transactionHistory)
	}
}

// This function will return a JSON response with a struct containing all relevant information
// about the balances in all bank accounts related to the user. If there is an error with the
// internet connection or database it will be reflected as an error response as well as a log
// to the server console
func GetBalance(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// create the user with the id obtained from middleware context
		userId := fmt.Sprintf("%v", r.Context().Value(models.Key("user"))) 

		// get all tokens with the user id
		tokens, err := models.GetTokens(app, userId)
		if err != nil {
			msg := "There was an error retrieving tokens from database affiliated with user"
			models.CreateError(w, http.StatusBadGateway, msg, err)
		}

		// define a struct for balance types
		type BType struct {
			Current		float64		`json:"current"`
			Name		string		`json:"name"`
			Limit		float64		`json:"limit,omitempty"`
		}

		// define a struct for balance totals
		type BTotal struct {
			Liquid		float64		`json:"liquid"`
			Credit		float64		`json:"credit"`
			Loan		float64		`json:"loan"`
			Total		float64		`json:"total"`
		}

		// define a struct for the balance response
		type Balance struct {
			Liquid		[]BType		`json:"liquid"`
			Credit		[]BType		`json:"credit"`
			Loan		[]BType		`json:"loan"`
			Net			BTotal		`json:"net"`
		}

		var balance Balance
		for i := range tokens {
			// get balance response
			res, err := app.Plaid.Client.GetBalances(tokens[i].Value)
			if err != nil {
				msg := "Error retrieving balance from client"
				models.CreateError(w, http.StatusBadGateway, msg, err)
			}
			
			// loop through all accounts related to that token
			for j := range res.Accounts {
				// determine what type it is and add it to the response
				switch(res.Accounts[j].Type) {
				case "depository": { 
					balance.Liquid = append(balance.Liquid, BType { 
						Current: 	res.Accounts[j].Balances.Current,
						Name:		res.Accounts[j].Name,
					})
					balance.Net.Liquid += res.Accounts[j].Balances.Current
					balance.Net.Total += res.Accounts[j].Balances.Current
				}  
				case "credit": {
					balance.Credit = append(balance.Credit, BType { 
						Current: 	res.Accounts[j].Balances.Current,
						Name:		res.Accounts[j].Name,
						Limit: 		res.Accounts[j].Balances.Limit,
					})  
					balance.Net.Credit += res.Accounts[j].Balances.Current
					balance.Net.Total -= res.Accounts[j].Balances.Current
				}
				case "loan": {
					balance.Loan = append(balance.Loan, BType { 
						Current: 	res.Accounts[j].Balances.Current,
						Name:		res.Accounts[j].Name,
					})  
					balance.Net.Loan += res.Accounts[j].Balances.Current
					balance.Net.Total -= res.Accounts[j].Balances.Current
				}
				}
			}
		}
	}
}