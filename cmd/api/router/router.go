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
	mux.GET("/v0/login", auth.Login(app))
	mux.GET("/v0/getLinkToken", m.Authenticate(p.GetPlaidToken(app), app))
	mux.GET("/v0/createAccessToken", m.Authenticate(p.CreateAccessToken(app), app))

	return mux
}