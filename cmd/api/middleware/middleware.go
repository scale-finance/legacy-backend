package middleware

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/julienschmidt/httprouter"
	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/cookie"
)

func Authenticate(next httprouter.Handle, app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if valid, err := cookie.Valid(r, app, "testcookie"); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Middleware",
				Message: "Failed to parse JWT",
			})
		} else if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Middleware",
				Message: "Invalid Authentication code",
			})
		} else {
			next(w,r,p)
		}
	}
}