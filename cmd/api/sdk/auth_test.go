package sdk_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/elopez00/scale-backend/cmd/api/sdk"
	"github.com/elopez00/scale-backend/pkg/test"

	"github.com/DATA-DOG/go-sqlmock"
)

func getBody() io.Reader {
	body, _ := json.Marshal(user)
	return bytes.NewBuffer(body)
}

func TestOnboardSuccess(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	// run expectation
	query := "INSERT INTO userinfo\\(id, firstname, lastname, email, password\\) VALUES\\(\\?,\\?,\\?,\\?,\\?\\)"
	app.DB.Mock.ExpectPrepare(query).ExpectExec()
	
	// create body of function
	body := getBody()
	res := test.Post("/onboard", sdk.Onboard(app), body)
	test.Response(t, res, http.StatusOK)
	test.MockExpectations(t, app)
}

func TestExistingUserError(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	// run expectation
	rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(user.Id, user.Email)
	
	query := `SELECT firstname, email FROM userinfo WHERE email \= \?`
	app.DB.Mock.ExpectQuery(query).WillReturnRows(rows)

	// create body of function
	body := getBody()
	res := test.Post("/onboard", sdk.Onboard(app), body)
	test.Response(t, res, http.StatusNotAcceptable)
	test.MockExpectations(t, app)
}

func TestPasswordIncorrect(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)
	
	password := "wrong password"

	// add rows to retrieve
	rows := sqlmock.NewRows([]string{"email", "password", "id"}).
		AddRow(user.Email, password, user.Id)
	
	// query
	query := `SELECT email, password, id FROM userinfo WHERE email \= \?`
	app.DB.Mock.ExpectQuery(query).WillReturnRows(rows)

	res := test.Post("/login", sdk.Login(app), getBody())
	test.Response(t, res, http.StatusUnauthorized)
	test.MockExpectations(t, app)
}

func TestUserInvalid(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	query := `SELECT email, password, id FROM userinfo WHERE email \= \?`
	app.DB.Mock.ExpectQuery(query)

	res := test.Post("/login", sdk.Login(app), getBody())
	test.Response(t, res, http.StatusNotFound)
	test.MockExpectations(t, app)
}

func TestUserSignout(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	res := test.GetWithCookie("/v0/logout", sdk.Logout(), app, "AuthToken")
	test.Response(t, res, http.StatusOK)
}

func TestUserSignoutFailure(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	res := test.Get("/v0/logout", sdk.Logout())
	test.Response(t, res, http.StatusBadRequest)
}