package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/application"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

// Authenticate is a function that takes in a handler and an application and returns another
// handler that tests the validity of authentication tokens found in auth cookie. If
// the token is expired or invalid for any reason the user will not be authenticated
// and will not be able to call api. Otherwise, the function will serve the res, req,
// parameters to the inputted handler and execute it.
func Authenticate(next httprouter.Handle, app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// checks if token is valid
		id, err := CookieIsValid(r, app, "AuthToken")
		if err != nil || len(id) == 0 {
			msg := "Unauthorized User"
			models.CreateError(w, http.StatusUnauthorized, msg, err)
			return
		}

		ctx := context.WithValue(r.Context(), models.Key("user"), id)
		r = r.WithContext(ctx)

		next(w, r, p)
	}
}

// CookieIsValid tests to see if cookie is still valid
func CookieIsValid(r *http.Request, app *application.App, name string) (string, error) {
	key := app.Config.GetServer()["key"]
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		err = errors.New("token invalid")
		return "", err
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	issuer := fmt.Sprint(claims["iss"])
	return issuer, nil
}
