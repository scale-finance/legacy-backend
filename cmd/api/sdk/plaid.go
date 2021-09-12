package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/application"

	"github.com/julienschmidt/httprouter"
	"github.com/plaid/plaid-go/plaid"
)

// GetPlaidToken returns the plaid token from authentication token. If in any case there is an error with
// the link token or the user's connection, it will return a json response error to the frontend
func GetPlaidToken(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// creates token configuration
		tokenConfig := plaid.LinkTokenConfigs{
			User: &plaid.LinkTokenUser{
				ClientUserID: GetIDFromContext(r),
			},
			ClientName:   "Scale",
			Products:     []string{"auth", "transactions"},
			CountryCodes: []string{app.Config.GetPlaid()["countryCode"]},
			Language:     "en",
			Webhook:      app.Plaid.RedirectURL,
		}

		// calls on the app's plaid client and creates a link token with the configuration
		// provided by the tokenConfig struct. If for whatever reason the client fails, it
		// will return a json response reflecting this issue
		tokenResponse, err := app.Plaid.Client.CreateLinkToken(tokenConfig)
		if err != nil {
			msg := "Failure to load client"
			models.CreateError(w, http.StatusBadGateway, msg, err)
			return
		}

		// when successful returns a result response
		msg := "Successfully received link token from plaid"
		models.CreateResponse(w, msg, tokenResponse.LinkToken)
	}
}

// UpdatePlaidToken returns a link token in update mode to allow user re-authentication
func UpdatePlaidToken(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer CloseBody(r)

		userId := GetIDFromContext(r)

		var token models.Token
		if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
			log.Println("Failed to decode token object from json")
			return
		}

		// get the token
		if err := token.Get(app, userId); err != nil {
			msg := "Failed to get token from database"
			models.CreateError(w, http.StatusNotFound, msg, err)
			return
		}

		tokenConfig := plaid.LinkTokenConfigs{
			User: &plaid.LinkTokenUser{
				ClientUserID: GetIDFromContext(r),
			},
			ClientName:   "Scale",
			CountryCodes: []string{app.Config.GetPlaid()["countryCode"]},
			Language:     "en",
			Webhook:      app.Plaid.RedirectURL,
			AccessToken: token.Value,
		}

		res, err := app.Plaid.Client.CreateLinkToken(tokenConfig)
		if err != nil {
			msg := "Failed to update token"
			models.CreateError(w, http.StatusBadGateway, msg, err)
			return
		}

		msg := "Successfully updated access token"
		models.CreateResponse(w, msg, res.LinkToken)
	}
}

// ExchangePublicToken this function takes care of creating the permanent access token
// that will be stored in the database for cross-platform connection to users' bank.
// If for whatever reason there is a problem with the client or public token, there
// are json responses and logs that will adequately reflect all issues
func ExchangePublicToken(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer CloseBody(r)

		var token models.Token
		err := json.NewDecoder(r.Body).Decode(&token)
		if err != nil {
			log.Println("Failed to decode token object from json")
			return
		}

		// Creates the permanent token using the public token it gets from the frontend
		// request body
		res, err := app.Plaid.Client.ExchangePublicToken(token.Value)
		if err != nil {
			msg := "Failure to exchange link token"
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

// GetTransactions is a function will get transactions from the past 12 months from all
// bank accounts affiliated with the user. If there is an error with the database retrieval
// or the plaid client call, this will be reflected in the json response accordingly.
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

		// initialize wait group
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
					log.Println(token.Value)
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

// GetBalance will return a JSON response with a struct containing all relevant information
// about the balances in all bank accounts related to the user. If there is an error with the
// internet connection or database it will be reflected as an error response as well as a log
// to the server console
func GetBalance(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// create the user with the id obtained from middleware context
		userId := GetIDFromContext(r)
		log.Println(userId)

		// get all tokens with the user id
		tokens, err := models.GetTokens(app, userId)
		if err != nil {
			msg := "There was an error retrieving tokens from database affiliated with user"
			models.CreateError(w, http.StatusBadGateway, msg, err)
		}
		
		// create waitGroup
		var waitGroup sync.WaitGroup

		// creating a context with cancel
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		// define a balance object and loop through tokens to start creating it
		var balance models.Balance

		// create continuous error
		var asyncError error
		var faultyInstitute string

		for _, token := range tokens {
			waitGroup.Add(1)

			go func(token *models.Token) {
				defer waitGroup.Done()

				// get balance response
				var backup plaid.GetAccountsResponse

				// stop calling requests if there has been an error somewhere
				select {
				case <-ctx.Done(): return
				default:
				}

				// get liabilities
				res, err := app.Plaid.Client.GetLiabilities(token.Value)

				// if getting liabilities fails because the product is not supported, then try to get 
				// the balances to at least have general info, and it will log it to the server. If the
				// liabilities call or the balances call fail for any other reason, it will return
				// as json response
				if err != nil && len(err.Error()) > 107 && err.Error()[85:107] == "PRODUCTS_NOT_SUPPORTED" {
					log.Println(err)
					if backup, err = app.Plaid.Client.GetAccounts(token.Value); err != nil {
						faultyInstitute = token.Id
						asyncError = err
						cancel()
						return
					}
				} else if err != nil {
					faultyInstitute = token.Id
					asyncError = err
					cancel()
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
		if ctx.Err() == nil {
			msg := "Successfully retrieved balance"
			models.CreateResponse(w, msg, balance)
		} else {
			msg := "Error retrieving Balances from client"
			models.CreateErrorWithResult(w, http.StatusBadGateway, msg, asyncError, faultyInstitute)
		}
	}
}
