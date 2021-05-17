package middleware 

import (
	"net/http"
	"fmt"

	"github.com/elopez00/scale-backend/models"
	jwt "github.com/dgrijalva/jwt-go"
)

func IsAuthorized(r *http.Request, env models.Env) bool {
	var output bool
	if cookie, err := r.Cookie("AuthToken"); err != nil {
		fmt.Println("Not Authorized")
		output = false
	} else {
		if token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(env.Key), nil
		}); err != nil { 
			fmt.Println("Invalid Signature")
			panic(err.Error())
		} else {
			output = token.Valid
		}
	}

	return output
}