package cookie

import (
	"time"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	application "github.com/elopez00/scale-backend/pkg/app"
)

// Creates valid httponly cookie
func Create(w http.ResponseWriter, app* application.App, name, id string) error {
	if token, err := generateJWT(app, id); err != nil {
		return err
	} else {
		cookie := http.Cookie {
			Name: name,
			Value: token,
			Expires: time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		}

		http.SetCookie(w, &cookie)
		return nil
	}
}

// tests to see if cookie is still valid
func Valid(r *http.Request, app *application.App, name string) (bool, error) {
	key := app.Config.GetKey()
	if cookie, err := r.Cookie(name); err != nil {
		return false, err
	} else {
		if token, err := jwt.ParseWithClaims(cookie.Value, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		}); err != nil {
			return false, err
		} else if !token.Valid {
			return false, err
		} else {
			return true, nil
		}
	}
} 

// static functions that generates JWT
func generateJWT(app *application.App, id string) (string, error) {
	key := app.Config.GetKey()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims {
		Issuer: 	id,
		ExpiresAt: 	time.Now().Add(time.Minute * 30).Unix(),
	})

	if token, err := claims.SignedString([]byte(key)); err != nil {
		return "", err
	} else {
		return token, nil
	}
}