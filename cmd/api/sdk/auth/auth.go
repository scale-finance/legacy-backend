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
		} else { 
			// finish gather important user data
			user.Password = EncryptPassword(user.Password)
			if len(user.Id) == 0 {
				user.Id = uuid.New().String()
			}
			// create user in database
			if err := user.Create(app); err != nil {
				// log.Println("Failed to create user in database: ", err)
				json.NewEncoder(w).Encode(models.Response {
					Status: 1,
					Type: "Onboarding",
					Message: "Unable to create user",
				}); return
			} else {
				json.NewEncoder(w).Encode(models.Response {
					Status: 0,
					Type: "Onboarding",
					Message: "User successfully created",
				})
			}
		}
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
		if err := authUser.GetCredentials(app, &actualUser); err != nil {
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Login",
				Message: "User Invalid",
			}); return
		} else { // check to see if passwords match
			if match := HashMatch(authUser.Password, actualUser.Password); !match {
				log.Println(actualUser.Password)
				json.NewEncoder(w).Encode(models.Response {
					Status: 1,
					Type: "Login",
					Message: "Password Incorrect",
				}); return
			}
		}

		// create a cookie to completely authenticate user
		if err := CreateCookie(w, app, "AuthToken", actualUser.Id); err != nil {
			log.Println(err.Error())
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Login",
				Message: "Failed to login",
			}); return
		} else {
			json.NewEncoder(w).Encode(models. Response {
				Status: 0,
				Type: "Login",
				Message: "User successfully authenticated",
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
			Expires: time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		}

		http.SetCookie(w, &cookie)
		return nil
	}
}

// static functions that generates JWT
func GenerateJWT(app *application.App, id string) (string, error) {
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