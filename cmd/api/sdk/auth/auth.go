package auth

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/julienschmidt/httprouter"
	application "github.com/elopez00/scale-backend/pkg/app"
	uuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Onboard(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.Exists(app) {
			json.NewEncoder(w).Encode(models.Response {
				Status: 1,
				Type: "Onboarding",
				Message: "User already exists",
			})
			return
		} else {
			if password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
				log.Println("Failed to generate password: ", err)
				json.NewEncoder(w).Encode(models.Response {
					Status: 1,
					Type: "Onboarding",
					Message: "Password encryption failure",
				})
				return
			} else {
				// finish gather important user data
				user.Id = uuid.New().String()
				user.Password = string(password)

				// create user in database
				if err := user.Create(app); err != nil {
					log.Println("Failed to creatue user in database: ", err)
					json.NewEncoder(w).Encode(models.Response {
						Status: 1,
						Type: "Onboarding",
						Message: "Unable to create user",
					})
					return
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
}