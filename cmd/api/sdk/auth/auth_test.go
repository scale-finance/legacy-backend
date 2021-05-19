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
		t.Errorf("handler returned wrong status code, what %v", http.StatusOK)
	} else {
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)
		
		// check if all sql sql expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	
		// check if the json response wasn't an error
		if (response.Status != 0) {
			t.Error("there was an error in authenticating user:", response.Message)
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
	if res := test.Post("/onboard", auth.Onboard(app), body); res.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code, got %v", res.Code)
	} else {
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)

		// check if all sql sql expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expections: %s", err)
		}

		// check if the json response wasn't an error
		if (response.Message != "User already exists") {
			t.Errorf("Expected %q got %q", "User already exists", response.Message)
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

	if res := test.Get("/login", auth.Login(app), getBody()); res.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code, got: %v", res.Code)
	} else {
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)

		// check if all sql sql expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expections: %s", err)
		}

		// check if the json response wasn't an error
		if (response.Message != "Password Incorrect") {
			t.Errorf("Expected %q got %q", "Password Incorrect", response.Message)
		}
	}
}

func TestUserInvalid(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query)

	if res := test.Get("/login", auth.Login(app), getBody()); res.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code, got: %v", res.Code)
	} else {
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)

		// check if all sql sql expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expections: %s", err)
		}

		// check if the json response wasn't an error
		if (response.Message != "User Invalid") {
			t.Errorf("Expected %q got %q", "User Invalid", response.Message)
		}
	}
}

// func TestCookie(t *testing.T) {
// 	app, mock := test.GetMockApp()
// 	defer app.DB.Client.Close()

// 	// tesing query
// 	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
// 	mock.ExpectQuery(query)

// 	if res := test.Get("/login", auth.Login(app), getBody()); res.Code != http.StatusOK {
// 		t.Errorf("Handler returned wrong status code, got: %v", res.Code)
// 	} else {
// 		if len(res.Result().Cookies()) == 0 {
// 			t.Error("Cookie was not created")
// 		}
// 	}
// }