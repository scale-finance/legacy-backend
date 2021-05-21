package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/elopez00/scale-backend/cmd/api/sdk/auth"
	"github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/cmd/api/sdk/plaid"
	"github.com/elopez00/scale-backend/cmd/api/middleware"
)
 
func Get(app *app.App) *httprouter.Router {
	mux := httprouter.New()
	mux.POST("/onboard", auth.Onboard(app))
	mux.GET("/login", auth.Login(app))
	mux.GET("/getLinkToken", middleware.Authenticate(plaid.GetPlaidToken(app), app))

	return mux
}