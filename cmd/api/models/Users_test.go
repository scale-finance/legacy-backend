package models_test

import (
	"testing"

	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/test"
	
	"github.com/DATA-DOG/go-sqlmock"
)

var user = models.User{
	Id:        "goingdowntosouthpark",
	FirstName: "Stan",
	LastName:  "Marsh",
	Email:     "smarsh@southpark.com",
	Password:  "southpark",
}

func TestUserCreate(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	query := "INSERT INTO userinfo\\(id, firstname, lastname, email, password\\) VALUES\\(\\?,\\?,\\?,\\?,\\?\\)"
	app.DB.Mock.ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, user.FirstName, user.LastName, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := user.Create(app)
	test.ModelMethod(t, err, "insert")
	test.MockExpectations(t, app)
}

func TestUserDoesExists(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(user.Id, user.Email)

	query := `SELECT firstname, email FROM userinfo WHERE email\="smarsh@southpark\.com"`
	app.DB.Mock.ExpectQuery(query).WillReturnRows(rows)

	exists := user.Exists(app)
	if !exists {
		t.Error("The user should exist")
		return
	}

	test.MockExpectations(t, app)
}

func TestUserDoesNotExist(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	query := `SELECT firstname, email FROM userinfo WHERE email\="smarsh@southpark\.com"`
	app.DB.Mock.ExpectQuery(query)

	exists := user.Exists(app)
	if exists {
		t.Error("User sould not exist")
		return
	}

	test.MockExpectations(t, app)
}

func TestGetCredential(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	rows := sqlmock.NewRows([]string{"email", "password", "id"}).
		AddRow(user.Email, user.Password, user.Id)

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	app.DB.Mock.ExpectQuery(query).WillReturnRows(rows)

	var actualUser models.User
	err := user.GetCredentials(app, &actualUser)
	test.ModelMethod(t, err, "select")
	if actualUser.Id != user.Id || actualUser.Password != user.Password {
		t.Error("Credentials are incorrect")
		return
	}
}

func TestUnsuccessfulGetCredential(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	app.DB.Mock.ExpectQuery(query)

	var actualUser models.User
	if err := user.GetCredentials(app, &actualUser); err == nil {
		t.Error("This process should have failed and returned an error")
		return
	}
}

