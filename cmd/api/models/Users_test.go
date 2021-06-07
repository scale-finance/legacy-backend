package models_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/test"
)

var user = models.User{
	Id:        "goingdowntosouthpark",
	FirstName: "Stan",
	LastName:  "Marsh",
	Email:     "smarsh@southpark.com",
	Password:  "southpark",
}

func TestUserCreate(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := "INSERT INTO userinfo\\(id, firstname, lastname, email, password\\) VALUES\\(\\?,\\?,\\?,\\?,\\?\\)"
	mock.ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, user.FirstName, user.LastName, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := user.Create(app)
	test.ModelMethod(t, err, "insert")
	test.MockExpectations(t, mock)
}

func TestUserDoesExists(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(user.Id, user.Email)

	query := `SELECT firstname, email FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	exists := user.Exists(app)
	if !exists {
		t.Error("The user should exist")
		return
	}

	test.MockExpectations(t, mock)
}

func TestUserDoesNotExist(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `SELECT firstname, email FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query)

	exists := user.Exists(app)
	if exists {
		t.Error("User sould not exist")
		return
	}

	test.MockExpectations(t, mock)
}

func TestGetCredential(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"email", "password", "id"}).
		AddRow(user.Email, user.Password, user.Id)

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	var actualUser models.User
	err := user.GetCredentials(app, &actualUser)
	test.ModelMethod(t, err, "select")
	if actualUser.Id != user.Id || actualUser.Password != user.Password {
		t.Error("Credentials are incorrect")
		return
	}
}

func TestUnsuccessfulGetCredential(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query)

	var actualUser models.User
	if err := user.GetCredentials(app, &actualUser); err == nil {
		t.Error("This process should have failed and returned an error")
		return
	}
}

