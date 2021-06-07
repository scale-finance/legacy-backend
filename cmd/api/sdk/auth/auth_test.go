package auth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/cmd/api/sdk/auth"
	"github.com/elopez00/scale-backend/pkg/test"
	
	"github.com/DATA-DOG/go-sqlmock"
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
	res := test.Post("/onboard", auth.Onboard(app), body)
	test.Response(t, res, http.StatusOK)
	test.MockExpectations(t, mock)
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
	res := test.Post("/onboard", auth.Onboard(app), body)
	test.Response(t, res, http.StatusNotAcceptable)
	test.MockExpectations(t, mock)
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

	res := test.Post("/login", auth.Login(app), getBody())
	test.Response(t, res, http.StatusUnauthorized)
	test.MockExpectations(t, mock)
}

func TestUserInvalid(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query)

	res := test.Post("/login", auth.Login(app), getBody())
	test.Response(t, res, http.StatusNotFound)
	test.MockExpectations(t, mock)
}

func TestUserSignout(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	res := test.GetWithCookie("/v0/logout", auth.Logout(app), app, "AuthToken") 
	test.Response(t, res, http.StatusOK)
}

func TestUserSignoutFailure(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	res := test.Get("/v0/logout", auth.Logout(app))
	test.Response(t, res, http.StatusBadRequest)
}