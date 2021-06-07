package auth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/cmd/api/sdk/auth"
	"github.com/elopez00/scale-backend/pkg/test"
)

var user = models.User {
	Id: "goingdowntosouthpark",
	FirstName: "Stan",
	LastName: "Marsh",
	Email: "smarsh@southpark.com",
	Password: "southpark",
}

func getBody() io.Reader {
	body, _ := json.Marshal(user)
	return bytes.NewBuffer(body)
}

func TestOnboardingSuccess(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()
	
	// run exepctation
	query := "INSERT INTO userinfo\\(id, firstname, lastname, email, password\\) VALUES\\(\\?,\\?,\\?,\\?,\\?\\)"
	mock.ExpectPrepare(query).ExpectExec()
	
	// create body of function
	body := getBody()
	if res := test.Post("/onboard", auth.Onboard(app), body); res.Code != http.StatusOK {
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("handler returned wrong status code, what %v", http.StatusOK)
		t.Error("there was an error in authenticating user:", response.Message)
	} else {
		// check if all sql sql expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}	
	}
}

func TestExistingUserError(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	// run expectation
	rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(user.Id, user.Email)
	
	query := `SELECT firstname, email FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	// create body of function
	body := getBody()
	if res := test.Post("/onboard", auth.Onboard(app), body); res.Code != http.StatusNotAcceptable {
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Handler returned wrong status code, got %v", res.Code)
		t.Errorf("Expected %q got %q", "User already exists", response.Message)
	} else {
		// check if all sql sql expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expections: %s", err)
		}
	}
}

func TestPasswordIncorrect(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()
	
	password := "wrong password"

	// add rows to retrieve
	rows := sqlmock.NewRows([]string{"email", "password", "id"}).
		AddRow(user.Email, password, user.Id)
	
	// query
	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	if res := test.Post("/login", auth.Login(app), getBody()); res.Code != http.StatusUnauthorized{
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Handler returned wrong status code, got: %v", res.Code)
		t.Errorf("Expected %q got %q", "Password Incorrect", response.Message)
	} else {
		// check if all sql sql expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expections: %s", err)
		}
	}
}

func TestUserInvalid(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query)

	if res := test.Post("/login", auth.Login(app), getBody()); res.Code != http.StatusNotFound {
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Handler returned wrong status code, got: %v", res.Code)
		t.Errorf("Expected %q got %q", "User Invalid", response.Message)
	} else {
		// check if all sql sql expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expections: %s", err)
		}
	}
}

func TestUserSignout(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	res := test.GetWithCookie("/v0/logout", auth.Logout(app), app, "AuthToken") 
	
	if res.Code != http.StatusOK {
		t.Error("This user was not successfully signed out")
	}
}

func TestUserSignoutFailure(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	res := test.Get("/v0/logout", auth.Logout(app))
	if res.Code != http.StatusBadRequest {
		t.Error("This function should not have successfully executed")
	}
}