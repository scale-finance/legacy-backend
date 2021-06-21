package plaid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"sync"

	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/application"

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

// This function will get transactions from the past 12 months from all bank accounts
// affiliated with the user. If there is an error with the database retrieval or the
// plaid client call, this will be reflected in the json response accordingly.
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
		startDate := time.Now().Add(-12 * 30 * 24 * time.Hour).Format(iso8601TimeFormat)
		endDate := time.Now().Format(iso8601TimeFormat)

		// intialize wait group 
		var waitGroup sync.WaitGroup

		// make transactions map that will be returned in result
		transactions := make(map[string][]plaid.Transaction)

		for _, token := range tokens {
			waitGroup.Add(1)

			go func(token *models.Token) {
				defer waitGroup.Done()

				// gets the transactions from plaid, if there is any error in the request,
				// it will be returned as a JSON response
				res, err := app.Plaid.Client.GetTransactions(token.Value, startDate, endDate)
				if err != nil {
					msg := "Failed to retrieve tokens from Plaid client"
					models.CreateError(w, http.StatusBadGateway, msg, err)
					return
				}
	
				// loop through all transactions and put them in transactions map
				for _, transaction := range res.Transactions {
					transactions[transaction.AccountID] = append(transactions[transaction.AccountID], transaction)
				}
			}(token)
		}
		waitGroup.Wait()

		msg := "Successfully retrieved transactions from all bank accounts"
		models.CreateResponse(w, msg, transactions)
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
		
		// create waitGroup
		var waitGroup sync.WaitGroup

		// define a balance object and loop through tokens to start creating it
		var balance models.Balance
		
		for _, token := range tokens {
			waitGroup.Add(1)
			
			go func(token *models.Token) {
				defer waitGroup.Done()

				// get balance response
				var backup plaid.GetAccountsResponse

				// get liabilities
				res, err := app.Plaid.Client.GetLiabilities(token.Value)

				// if getting liabilities fails because the product is not supported, then try to get 
				// the balances to at least have general info and it will log it to the server. If the 
				// liabilities call or the balances call fail for any other reason, it will return
				// as json response
				if err != nil && len(err.Error()) > 107 && err.Error()[85:107] == "PRODUCTS_NOT_SUPPORTED" {
					log.Println(err) 
					
					if backup, err = app.Plaid.Client.GetAccounts(token.Value); err != nil {
						msg := "Error retrieving Liabilities from client:"
						models.CreateError(w, http.StatusBadGateway, msg, err)
						return
					}
				} else if err != nil { 
					msg := "Error retrieving Liabilities from client"
					models.CreateError(w, http.StatusBadGateway, msg, err)
					return
				}

				// loop through all accounts related to that token
				if len(backup.Accounts) == 0 {
					liabilities := models.PlaidLiabilities(res.Liabilities)

					for _, account := range res.Accounts {
						balance.AddBalance(token.Institution, account, &liabilities)
					}
				}

				for _, account := range res.Accounts {
					balance.AddBalance(token.Institution, account, &models.PlaidLiabilities{})
				}
			}(token)
		}

		waitGroup.Wait()

		// return balance object in result of response
		msg := "Successfully retrieved balance"
		models.CreateResponse(w, msg, balance)
	}
}
