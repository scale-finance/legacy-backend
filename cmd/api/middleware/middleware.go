package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/elopez00/scale-backend/cmd/api/models"
	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/pkg/cookie"
	"github.com/julienschmidt/httprouter"
)

func Authenticate(next httprouter.Handle, app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		if id, err := cookie.Valid(r, app, "AuthToken"); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Middleware",
				Message: "Unauthorized User",
			}); return
		} else if len(id) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Middleware",
				Message: "Unauthorized User",
			}); return
		} else {			
			ctx := context.WithValue(r.Context(), models.Key("user"), id)
			r = r.WithContext(ctx)
			
			next(w, r, p)
		}
	}
}