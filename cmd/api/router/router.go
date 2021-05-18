package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/elopez00/scale-backend/cmd/api/sdk/auth"
	"github.com/elopez00/scale-backend/pkg/app"
)
 
func Get(app *app.App) *httprouter.Router {
	mux := httprouter.New()
	mux.POST("/onboard", auth.Onboard(app))

	return mux
}