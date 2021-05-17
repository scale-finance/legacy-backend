package database_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"bytes"
	"log"
	"os"

	"github.com/elopez00/scale-backend/database"
	"github.com/elopez00/scale-backend/helper"
	"github.com/elopez00/scale-backend/models"
	"github.com/gorilla/mux"
)

type App struct {
	DB		*sql.DB
	Router	*mux.Router
}

var app App
var env models.Env

func TestMain(m *testing.M) {
	env = helper.InitTestEnv()

	app.Router 	= mux.NewRouter()
	app.DB		= database.Connect(app.Router, env)

	code := m.Run()
	os.Exit(code)
}

func TestDBConnection(t *testing.T) { 
	query := "SELECT * FROM userinfo"
	if _, err := app.DB.Query(query); err != nil {
		t.Fatal("DB authentication failure")
		log.Fatal(err)
	}
}

func TestOnboarding(t *testing.T) {
	testUser, _ := json.Marshal(models.User {
		Email: "smarshal@gmail.com",
		FirstName: "Stan",
		LastName: "Marshal",
		Password: "southpark",
	})

	req, _ := http.NewRequest("POST", "/onboard", bytes.NewBuffer(testUser))
	res := helper.ExecuteRequest(req, app.Router)

	helper.CheckResponse(t, http.StatusOK, res.Code)
}

func TestPasswordEncrypt(t *testing.T) {
	query := "SELECT password FROM userinfo WHERE firstname='Stan';"
	if result, err := app.DB.Query(query); err != nil {
		panic(err.Error())
	} else {
		var user models.User
		for result.Next() {
			if err := result.Scan(&user.Password); err != nil {
				panic(err.Error())
			} else {
				if user.Password == "southpark" {
					t.Error("Password was not encrypted successfully")
				}
			}
		}
	}
}