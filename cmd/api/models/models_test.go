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

func TestUpdateWhiteList(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := 
		`INSERT INTO whitelist\(id, category, name, itemId\) VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ON ` +
		`DUPLICATE KEY UPDATE id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	
	budget := models.Budget {
		Request: models.UpdateRequest {
			Update: models.UpdateObject {
				WhiteList: []models.WhiteListItem {
					{ whitelist.Category, whitelist.Name, "hellobro" },
					{ whitelist.Category, whitelist.Name, "goodbyeo" },
				},
			},
		},
	}

	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(
			user.Id, budget.Request.Update.WhiteList[0].Category, budget.Request.Update.WhiteList[0].Name, budget.Request.Update.WhiteList[0].Id, 
			user.Id, budget.Request.Update.WhiteList[1].Category, budget.Request.Update.WhiteList[1].Name, budget.Request.Update.WhiteList[1].Id, 
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := budget.Update(app, user.Id); err != nil {
		t.Error("Error inserting information to database:", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("There were unfulfilled expectations:", err)
		return
	}
}

func TestUpdateCategory(t *testing.T) {
	app, mock := test.GetMockApp() 
	defer app.DB.Client.Close()
	
	query := 
		`INSERT INTO categories\(id, name, budget, categoryId\) VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ON ` +
		`DUPLICATE KEY UPDATE id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId;`
	
	budget := models.Budget {
		Request: models.UpdateRequest {
			Update: models.UpdateObject {
				Categories: []models.Category {
					{ whitelist.Category, 100, "catie" },
					{ "shopping", 200, "cattegorcatie" },
				},
			},
		},
	}
	
	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(
			user.Id, budget.Request.Update.Categories[0].Name, budget.Request.Update.Categories[0].Budget, budget.Request.Update.Categories[0].Id,
			user.Id, budget.Request.Update.Categories[1].Name, budget.Request.Update.Categories[1].Budget, budget.Request.Update.Categories[1].Id,	
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := budget.Update(app, user.Id); err != nil {
		t.Error("Error inserting data into data:", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("There were unfulfilled expectations:", err)
		return
	}
}

func TestUpdateBudget(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	budget := models.Budget {
		Request: models.UpdateRequest {
			Update: models.UpdateObject {
				Categories: []models.Category { 
					{ Name: "Shopping", Budget: 300, Id: "a;sldfkdj" },
					{ Name: "Fast Food", Budget: 100, Id: "asldfkjs" },
				},
				WhiteList: []models.WhiteListItem {
					{ Name: "Polo Store", Category: "Shopping", Id: ";lkj3lk" },
					{ Name: "Five Guys", Category: "Fast Food", Id: ";lkj;fl" },
					{ Name: "Chipotle", Category: "Fast Food", Id: "a;sdlf6k" },
				},
			},
		},
	}
	
	query2 := 
		`INSERT INTO categories\(id, name, budget, categoryId\) VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ON ` +
		`DUPLICATE KEY UPDATE id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId;`
	mock.
		ExpectPrepare(query2).
		ExpectExec().
		WithArgs(
			user.Id, budget.Request.Update.Categories[0].Name, budget.Request.Update.Categories[0].Budget, budget.Request.Update.Categories[0].Id,
			user.Id, budget.Request.Update.Categories[1].Name, budget.Request.Update.Categories[1].Budget, budget.Request.Update.Categories[1].Id,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	query1 := 
		`INSERT INTO whitelist\(id, category, name, itemId\) VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ON ` +
		`DUPLICATE KEY UPDATE id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	mock.
		ExpectPrepare(query1).
		ExpectExec().
		WithArgs(
			user.Id, budget.Request.Update.WhiteList[0].Category, budget.Request.Update.WhiteList[0].Name, budget.Request.Update.WhiteList[0].Id,
			user.Id, budget.Request.Update.WhiteList[1].Category, budget.Request.Update.WhiteList[1].Name, budget.Request.Update.WhiteList[1].Id,
			user.Id, budget.Request.Update.WhiteList[2].Category, budget.Request.Update.WhiteList[2].Name, budget.Request.Update.WhiteList[2].Id,
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