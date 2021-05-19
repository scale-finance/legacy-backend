package models_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elopez00/scale-backend/cmd/api/models"
	application "github.com/elopez00/scale-backend/pkg/app"
)

var user = models.User {
	Id: "goingdowntosouthpark",
	FirstName: "Stan",
	LastName: "Marsh",
	Email: "smarsh@southpark.com",
	Password: "southpark",
}

func newMock() (*application.App, sqlmock.Sqlmock){
	if db, mock, err := sqlmock.New(); err != nil {
		panic(err.Error)
	} else {
		app := application.GetTest(db)
		return app, mock
	}
}

func TestUserCreate(t *testing.T) {
	app, mock := newMock()
	defer app.DB.Client.Close()

	query := "INSERT INTO userinfo\\(id, firstname, lastname, email, password\\) VALUES\\(\\?,\\?,\\?,\\?,\\?\\)"
	mock.ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, user.FirstName, user.LastName, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(0,0))
	
	
	user.Create(app)
	if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expections: %s", err)
    }
}

func TestUserDoesExists(t *testing.T) {
	app, mock := newMock()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"id","email"}).
		AddRow(user.Id, user.Email)

	query := `SELECT firstname, email FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	exists := user.Exists(app)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	} else if !exists {
		t.Error("The user should exist")
	}
}

func TestUserDoesNotExist(t *testing.T) {
	app, mock := newMock()
	defer app.DB.Client.Close()

	query := `SELECT firstname, email FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query)

	exists := user.Exists(app)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	} else if exists {
		t.Error("User sould not exist")
	}
}

func TestSuccessfulGetCredential(t *testing.T) {
	app, mock := newMock()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"email", "password", "id"}).
		AddRow(user.Email, user.Password, user.Id)

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	var actualUser models.User
	if err := user.GetCredentials(app, &actualUser); err != nil {
		t.Errorf("Process should have run without errors, instead got: %v", err.Error())
	} else if actualUser.Id != user.Id || actualUser.Password != user.Password {
		t.Error("Credentials are incorrect")
	}
}

func TestUnsuccessfulGetCredential(t *testing.T) {
	app, mock := newMock()
	defer app.DB.Client.Close()

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query)

	var actualUser models.User
	if err := user.GetCredentials(app, &actualUser); err == nil {
		t.Error("This process should have failed and returned an error")
	}
}