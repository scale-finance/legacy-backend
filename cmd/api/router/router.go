package router

import (
	m "github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/sdk"
	"github.com/elopez00/scale-backend/pkg/application"
	
	"github.com/julienschmidt/httprouter"
)

// Get retrieves the api routes used in server
func Get(app *application.App) *httprouter.Router {
	mux := httprouter.New()

	// registration and account management
	mux.POST("/v0/onboard", sdk.Onboard(app))
	mux.POST("/v0/login", sdk.Login(app))
	mux.GET("/v0/logout", sdk.Logout())

	// plaid token management
	mux.POST("/v0/token/exchange", m.Authenticate(sdk.ExchangePublicToken(app), app))
	mux.GET("/v0/token/link", m.Authenticate(sdk.GetPlaidToken(app), app))

	// transactions
	mux.GET("/v0/transactions", m.Authenticate(sdk.GetTransactions(app), app))

	// balances
	mux.GET("/v0/balances", m.Authenticate(sdk.GetBalance(app), app))

	// budget
	mux.GET("/v0/budget", m.Authenticate(sdk.GetBudget(app), app))
	mux.PUT("/v0/budget", m.Authenticate(sdk.UpdateBudget(app), app))
	mux.DELETE("/v0/budget", m.Authenticate(sdk.UpdateBudget(app), app))

	// temp
	mux.GET("/v0/", m.Authenticate(sdk.AuthCheck(), app))

	return mux
}