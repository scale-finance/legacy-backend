package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/cmd/api/sdk/auth"
	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/pkg/test"
	// "github.com/julienschmidt/httprouter"
)

var user models.User

func TestOnboardingSuccess(t *testing.T) {
	db, mock, _ := sqlmock.New()
	app := application.GetTest(db)
	defer app.DB.Client.Close()

	// instantiate user
	user = models.User {
		FirstName: "Stan",
		LastName: "Marsh",
		Email: "smarsh@southpark.com",
		Password: "southpark",
	}
	
	// run exepctation
	query := "INSERT INTO userinfo\\(id, firstname, lastname, email, password\\) VALUES\\(\\?,\\?,\\?,\\?,\\?\\)"
	mock.ExpectPrepare(query).ExpectExec()
	
	// create body of function
	body, _ := json.Marshal(user)
	res := test.Post("/onboard", auth.Onboard(app), bytes.NewBuffer(body))

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code, got %v want %v", status, http.StatusOK)
	}

	var response models.Response 
	json.NewDecoder(res.Body).Decode(&response)

	if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expections: %s", err)
    }

	if (response.Status != 0) {
		t.Error("there was an error in authenticating user")
	}
}