package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/elopez00/scale-backend/cmd/api/sdk/auth"
	"github.com/elopez00/scale-backend/pkg/app"
	p "github.com/elopez00/scale-backend/cmd/api/sdk/plaid"
	m "github.com/elopez00/scale-backend/cmd/api/middleware"
)

// Gets api routes used in server
func Get(app *app.App) *httprouter.Router {
	mux := httprouter.New()
	mux.POST("/v0/onboard", auth.Onboard(app))
	mux.POST("/v0/login", auth.Login(app))
	mux.POST("/v0/exchangeAccessToken", m.Authenticate(p.CreateAccessToken(app), app))
	mux.GET("/v0/getLinkToken", m.Authenticate(p.GetPlaidToken(app), app))
	mux.GET("/v0/getTransactions", m.Authenticate(p.GetTransactions(app), app))
	mux.GET("/v0/logout", auth.Logout(app))
	mux.GET("/v0/", m.Authenticate(auth.AuthCheck(), app))

	return mux
}