package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"log"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/elopez00/scale-backend/cmd/api/models"
	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/julienschmidt/httprouter"
)

// TODO use alice or find a way to chain desired handlers to be affected by middleware

// Middleware function that takes in a handler and an application and returns another
// handler that tests the validity of authentication tokens found in auth cookie. If
// the token is expired or invalid for any reason the user will not be authenticated 
// and will not be able to call api. Otherwise, the function will serve the res, req,
// parameters to the inputted handler and execute it.
func Authenticate(next httprouter.Handle, app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		// checks if token is valid
		if id, err := CookieIsValid(r, app, "AuthToken"); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Error retrieving cookie:", err.Error())
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Middleware",
				Message: "Unauthorized User",
			}); return 
		} else if len(id) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("No access token found")
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


// tests to see if cookie is still valid
func CookieIsValid(r *http.Request, app *application.App, name string) (string, error) {
	key := app.Config.GetKey()
	if cookie, err := r.Cookie(name); err != nil {
		return "", err
	} else {
		if token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		}); err != nil {
			return "", err
		} else if !token.Valid {
			return "", err
		} else {
			claims, _ := token.Claims.(jwt.MapClaims)
			issuer := fmt.Sprint(claims["iss"])
			return issuer, nil
		}
	}
}