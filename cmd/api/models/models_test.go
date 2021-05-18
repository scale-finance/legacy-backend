package models_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elopez00/scale-backend/cmd/api/models"
	application "github.com/elopez00/scale-backend/pkg/app"
)


func TestCreate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	app := application.GetTest(db)
	fmt.Println(app.DB.Client)

	defer app.DB.Client.Close()

	user := models.User {
		FirstName: "Stan",
		LastName: "Marsh",
		Email: "smarsh@southpark.com",
		Password: "southpark",
	}
	
	query := "INSERT INTO userinfo\\(id, firstname, lastname, email, password\\) VALUES\\(\\?,\\?,\\?,\\?,\\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.
	ExpectExec().
	WithArgs(user.Id, user.FirstName, user.LastName, user.Email, user.Password).
	WillReturnResult(sqlmock.NewResult(0,0))
	
	
	user.Create(app)
	if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expections: %s", err)
    }
}