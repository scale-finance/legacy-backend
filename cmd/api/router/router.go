package router

import (
	m "github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/sdk"
	"github.com/elopez00/scale-backend/pkg/application"
	
	"github.com/julienschmidt/httprouter"
)

// Gets api routes used in server
func Get(app *application.App) *httprouter.Router {
	mux := httprouter.New()
	mux.POST("/v0/onboard", sdk.Onboard(app))
	mux.POST("/v0/login", sdk.Login(app))
	mux.POST("/v0/token/exchange", m.Authenticate(sdk.ExchangePublicToken(app), app))
	mux.POST("/v0/budget/update", m.Authenticate(sdk.UpdateBudget(app), app))
	mux.GET("/v0/token/link", m.Authenticate(sdk.GetPlaidToken(app), app))
	mux.GET("/v0/transactions/get", m.Authenticate(sdk.GetTransactions(app), app))
	mux.GET("/v0/balances/get", m.Authenticate(sdk.GetBalance(app), app))
	mux.GET("/v0/budget/get", m.Authenticate(sdk.GetBudget(app), app))
	mux.GET("/v0/logout", sdk.Logout(app))
	mux.GET("/v0/", m.Authenticate(sdk.AuthCheck(), app))

	return mux
}
