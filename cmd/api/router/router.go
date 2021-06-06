package router

import (
	m "github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/sdk/auth"
	b "github.com/elopez00/scale-backend/cmd/api/sdk/budgeting"
	p "github.com/elopez00/scale-backend/cmd/api/sdk/plaid"
	"github.com/elopez00/scale-backend/pkg/app"
	"github.com/julienschmidt/httprouter"
)

// Gets api routes used in server
func Get(app *app.App) *httprouter.Router {
	mux := httprouter.New()
	mux.POST("/v0/onboard", auth.Onboard(app))
	mux.POST("/v0/login", auth.Login(app))
	mux.POST("/v0/exchangePublicToken", m.Authenticate(p.ExchangePublicToken(app), app))
	mux.POST("/v0/updateBudget", m.Authenticate(b.Update(app), app))
	mux.GET("/v0/getLinkToken", m.Authenticate(p.GetPlaidToken(app), app))
	mux.GET("/v0/getTransactions", m.Authenticate(p.GetTransactions(app), app))
	mux.GET("/v0/getBalances", m.Authenticate(p.GetBalance(app), app))
	mux.GET("/v0/logout", auth.Logout(app))
	mux.GET("/v0/", m.Authenticate(auth.AuthCheck(), app))

	return mux
}
