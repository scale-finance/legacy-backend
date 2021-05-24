package auth

import (
	"encoding/json"
	"net/http"
	"time"
	"log"
	
	application "github.com/elopez00/scale-backend/pkg/app"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/julienschmidt/httprouter"
	"github.com/google/uuid"
)

// onboard user to DB given application sequence. This function is in charge of creating
// a new user in the database (given one does not already exist with the same credentials)
// and will give each user a unique ID and a hashed password for further authentication
func Onboard(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		
		// get user input from body
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check to see if user already exists in database
		if user.Exists(app) {
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Onboarding",
				Message: "User already exists",
			}); return
		}

		// finish gather important user data
		user.Password = EncryptPassword(user.Password)
		if len(user.Id) == 0 {
			user.Id = uuid.New().String()
		}

		// create user in database
		err := user.Create(app)
		if err != nil {
			// log.Println("Failed to create user in database: ", err)
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Onboarding",
				Message: "Unable to create user",
			})
			return
		} 

		// create a cookie to completely authenticate user
		err = CreateCookie(w, app, "AuthToken", user.Id)
		if err != nil {
			log.Println(err.Error())
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Login",
				Message: "Failed to login",
			})
			return
		}

		// return successful onboarding message
		json.NewEncoder(w).Encode(models.Response {
			Status: 0,
			Type: "Onboarding",
			Message: "User successfully created",
		})
	}
}

// logs in user by doing a preliminary check to backend to check if the user exists. After
// verification, the function will compare hashed and input password so it can then focus
// on creating a jwt token
func Login(app *application.App) httprouter.Handle {
	return func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		// grabs input user from body
		var authUser models.User
		if err := json.NewDecoder(r.Body).Decode(&authUser); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// get actual user from database
		var actualUser models.User
		err := authUser.GetCredentials(app, &actualUser)
		if err != nil {
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Login",
				Message: "User Invalid",
			})
			return
		}

		// check to see if passwords match
		if match := HashMatch(authUser.Password, actualUser.Password); !match {
			log.Println(actualUser.Password)
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Login",
				Message: "Password Incorrect",
			})
			return
		}

		// create a cookie to completely authenticate user
		err = CreateCookie(w, app, "AuthToken", actualUser.Id)
		if err != nil {
			log.Println(err.Error())
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Login",
				Message: "Failed to login",
			})
			return
		} 

		// send successful authentication message to client
		json.NewEncoder(w).Encode(models. Response {
			Status: 0,
			Type: "Login",
			Message: "User successfully authenticated",
		})
	}
}

// This function logs the user out of their session by deleting the AuthToken cookie containing
// the user user's JWT. The function should return an error status if there is no token to delete
func Logout(app *application.App) httprouter.Handle {
	return func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if len(r.Cookies()) < 1 {
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Sign out",
				Message: "User already signed out",
			})
		} else {
			DeleteCookie(w, app, "AuthToken")
			json.NewEncoder(w).Encode(models.Response {
				Status: 0,
				Type: "Sign out",
				Message: "User successfully signed out",
			})
		}
	}
}

// Creates valid httponly cookie
func CreateCookie(w http.ResponseWriter, app* application.App, name, id string) error {
	if token, err := GenerateJWT(app, id); err != nil {
		return err
	} else {
		cookie := http.Cookie {
			Name: name,
			Value: token,
			Expires: time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		}

		http.SetCookie(w, &cookie)
		return nil
	}
}

// Deletes existing cookie
func DeleteCookie(w http.ResponseWriter, app* application.App, name string) {
	cookie := http.Cookie {
		Name: name,
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
}

// Function that tests authentication state
func AuthCheck() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		json.NewEncoder(w).Encode(models.Response {
			Status: 0,
			Type: "Authentication",
			Message: "This app is authenticated",
		})
	}
}

// This function generates a JWT token based on the user id stored in the database that expires in
// 24 hours. If this function were to fail, its error would be returned respectfully.
func GenerateJWT(app *application.App, id string) (string, error) {
	key := app.Config.GetKey()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims {
		Issuer: 	id,
		ExpiresAt: 	time.Now().Add(24 * time.Hour).Unix(),
	})

	if token, err := claims.SignedString([]byte(key)); err != nil {
		return "", err
	} else {
		return token, nil
	}
}