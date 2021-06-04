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

var token = models.Token{
	Value: "randomaccess",
	Id:    "randomid",
}

var whitelist = models.WhiteListItem{
	Category: "Shopping",
	Name:     "Source 1",
}

// * user tests

func TestUserCreate(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := "INSERT INTO userinfo\\(id, firstname, lastname, email, password\\) VALUES\\(\\?,\\?,\\?,\\?,\\?\\)"
	mock.ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, user.FirstName, user.LastName, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(0, 0))

	user.Create(app)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
		return
	}
}

func TestUserDoesExists(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(user.Id, user.Email)

	query := `SELECT firstname, email FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	exists := user.Exists(app)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !exists {
		t.Error("The user should exist")
		return
	}
}

func TestUserDoesNotExist(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `SELECT firstname, email FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query)

	exists := user.Exists(app)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
		return
	}

	if exists {
		t.Error("User sould not exist")
		return
	}
}

func TestSuccessfulGetCredential(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"email", "password", "id"}).
		AddRow(user.Email, user.Password, user.Id)

	query := `SELECT email, password, id FROM userinfo WHERE email\="smarsh@southpark\.com"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	var actualUser models.User
	if err := user.GetCredentials(app, &actualUser); err != nil {
		t.Errorf("Process should have run without errors, instead got: %v", err.Error())
		return
	}

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

// * token tests

func TestTokenAdd(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `INSERT INTO plaidtokens\(id, token, itemID\) VALUES\(\?,\?,\?\)`
	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, token.Value, token.Id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	token.Add(app, user.Id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
		return
	}
}

func TestGetTokens(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"id", "token", "itemID"}).
		AddRow(user.Id, "token1", "id1").
		AddRow(user.Id, "token2", "id2")

	query := `SELECT id, token, itemID FROM plaidtokens WHERE id\="goingdowntosouthpark"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	tokens, _ := models.GetTokens(app, user.Id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
		return
	}

	if len(tokens) == 0 {
		t.Error("The function did not return the tokens")
		return
	}
}

// * budget tests

func TestAddWhiteList(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `INSERT INTO whitelist\(id, category, name\) VALUES \(\?,\?,\?\), \(\?,\?,\?\)`
	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, whitelist.Category, whitelist.Name, user.Id, whitelist.Category, whitelist.Name).
		WillReturnResult(sqlmock.NewResult(0, 0))

	list := []models.WhiteListItem {
		{ whitelist.Category, whitelist.Name, },
		{ whitelist.Category, whitelist.Name, },
	}

	if err := models.AddWhiteList(app, user.Id, list); err != nil {
		t.Error("Error inserting information to database:", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("There were unfulfilled expectations:", err)
		return
	}
}

func TestCreateBudget(t *testing.T) {
	app, mock := test.GetMockApp() 
	defer app.DB.Client.Close()
	
	query := `INSERT INTO categories\(id, name, budget\) VALUES \(\?,\?,\?\), \(\?,\?,\?\)`
	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, whitelist.Category, float64(100), user.Id, "shopping", float64(200)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	
	budget := models.Budget {
		Categories: []models.Category {
			{ whitelist.Category, float64(100), },
			{ "shopping", float64(200), },
		},
	}

	if err := budget.Create(app, user.Id); err != nil {
		t.Error("Error inserting data into data:", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("There were unfulfilled expectations:", err)
		return
	}
}

func TestUpdateBudgetAdd(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	budget := models.Budget {
		Request: models.UpdateRequest {
			Add: models.UpdateObject {
				Categories: []models.UpdateCategory { 
					{ New: models.Category { Name: "Shopping", Budget: 300 }, },
					{ New: models.Category { Name: "Fast Food", Budget: 100 }, },
				},
				WhiteList: []models.UpdateWhiteList {
					{ New: models.WhiteListItem { Name: "Polo Store", Category: "Shopping"}, },
					{ New: models.WhiteListItem { Name: "Five Guys", Category: "Fast Food" }, },
					{ New: models.WhiteListItem { Name: "Chipotle", Category: "Fast Food" }, },
				},
			},
		},
	}
	
	query2 := `INSERT INTO categories\(id, name, budget\) VALUES \(\?,\?,\?\), \(\?,\?,\?\)`
	mock.
		ExpectPrepare(query2).
		ExpectExec().
		WithArgs(
			user.Id, budget.Request.Add.Categories[0].New.Name, budget.Request.Add.Categories[0].New.Budget,
			user.Id, budget.Request.Add.Categories[1].New.Name, budget.Request.Add.Categories[1].New.Budget,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	query1 := `INSERT INTO whitelist\(id, category, name\) VALUES \(\?,\?,\?\), \(\?,\?,\?\), \(\?,\?,\?\)`
	mock.
		ExpectPrepare(query1).
		ExpectExec().
		WithArgs(
			user.Id, budget.Request.Add.WhiteList[0].New.Category, budget.Request.Add.WhiteList[0].New.Name,
			user.Id, budget.Request.Add.WhiteList[1].New.Category, budget.Request.Add.WhiteList[1].New.Name,
			user.Id, budget.Request.Add.WhiteList[2].New.Category, budget.Request.Add.WhiteList[2].New.Name,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := budget.Update(app, user.Id); err != nil {
		t.Error("There was an error updating the budget:", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("There were unfulfilled expectations:", err)
		return
	}
}