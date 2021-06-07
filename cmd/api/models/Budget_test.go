package models_test

import (
	"testing"

	"github.com/elopez00/scale-backend/pkg/test"
	"github.com/elopez00/scale-backend/cmd/api/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUpdateWhiteList(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := 
		`INSERT INTO whitelist\(id, category, name, itemId\) ` + 
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	
	budget := models.Budget {
		Request: models.UpdateRequest {
			Update: models.UpdateObject {
				WhiteList: []models.WhiteListItem {
					{ "Shopping", "Calvin Klien", "hellobro" },
					{ "Taxes", "IRS", "goodbyeo" },
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

	err := budget.Update(app, user.Id)
	test.ModelMethod(t, err, "insert")
	test.MockExpectations(t, mock)
}

func TestUpdateCategory(t *testing.T) {
	app, mock := test.GetMockApp() 
	defer app.DB.Client.Close()
	
	query := 
		`INSERT INTO categories\(id, name, budget, categoryId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId;`
	
	budget := models.Budget {
		Request: models.UpdateRequest {
			Update: models.UpdateObject {
				Categories: []models.Category {
					{ "Taxes", 100, "catie" },
					{ "Shopping", 200, "cattegorcatie" },
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

	err := budget.Update(app, user.Id)
	test.ModelMethod(t, err, "insert")
	test.MockExpectations(t, mock)
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
		`INSERT INTO categories\(id, name, budget, categoryId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId;`
	mock.
		ExpectPrepare(query2).
		ExpectExec().
		WithArgs(
			user.Id, budget.Request.Update.Categories[0].Name, budget.Request.Update.Categories[0].Budget, budget.Request.Update.Categories[0].Id,
			user.Id, budget.Request.Update.Categories[1].Name, budget.Request.Update.Categories[1].Budget, budget.Request.Update.Categories[1].Id,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	query1 := 
		`INSERT INTO whitelist\(id, category, name, itemId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	mock.
		ExpectPrepare(query1).
		ExpectExec().
		WithArgs(
			user.Id, budget.Request.Update.WhiteList[0].Category, budget.Request.Update.WhiteList[0].Name, budget.Request.Update.WhiteList[0].Id,
			user.Id, budget.Request.Update.WhiteList[1].Category, budget.Request.Update.WhiteList[1].Name, budget.Request.Update.WhiteList[1].Id,
			user.Id, budget.Request.Update.WhiteList[2].Category, budget.Request.Update.WhiteList[2].Name, budget.Request.Update.WhiteList[2].Id,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := budget.Update(app, user.Id)
	test.ModelMethod(t, err, "insert")
	test.MockExpectations(t, mock)
}

func TestDeleteCategoryAndListItems(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	budget := models.Budget {
		Request: models.UpdateRequest {  
			Remove: models.UpdateObject {
				Categories: []models.Category {
					{ "Shopping", 400, "cid123" },
				},
				WhiteList: []models.WhiteListItem {
					{ "cid123", "Calvin Klein", "something" },
					{ "cid123", "Ralph Lauren", "domething" },
				},
			},
		},
	}

	query := 
		`DELETE categories, whitelist ` +
		`FROM categories LEFT JOIN whitelist ` +
		`ON categories\.categoryId \= whitelist\.category ` +
		`WHERE categories\.id \= \? AND categories\.categoryId IN \(\?\);`

	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, budget.Request.Remove.Categories[0].Id).
		WillReturnResult(sqlmock.NewResult(0, 3))
	
	err := models.Delete(app, user.Id, budget)
	test.ModelMethod(t, err, "delete")
	test.MockExpectations(t, mock)
}

func TestDeleteWhiteListItem(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	budget := models.Budget {
		Request: models.UpdateRequest {
			Remove: models.UpdateObject {
				WhiteList: []models.WhiteListItem {
					{ "cid123", "Calvin Klein", "something" },
					{ "cid123", "Ralph Lauren", "domething" },
				},
			},
		},
	}

	query := `DELETE FROM whitelist WHERE whitelist\.id \= \? AND whitelist\.itemId IN \(\?,\?\);`
	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(
			user.Id,
			budget.Request.Remove.WhiteList[0].Id,
			budget.Request.Remove.WhiteList[1].Id,
		).
		WillReturnResult(sqlmock.NewResult(0, 2))
	
	err := models.Delete(app, user.Id, budget)
	test.ModelMethod(t, err, "delete")
	test.MockExpectations(t, mock)
}

func TestDeleteWhiteListAndCategories(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	budget := models.Budget {
		Request: models.UpdateRequest {
			Remove: models.UpdateObject {
				Categories: []models.Category {
					{ "Shopping", 400, "cid123" },
				},
				WhiteList: []models.WhiteListItem {
					{ "cid123", "Calvin Klein", "something" },
					{ "cid123", "Ralph Lauren", "domething" },
					{ "cid456", "Calvin Klein", "comething" },
				},
			},
		},
	}

	query1 := 
		`DELETE categories, whitelist ` +
		`FROM categories LEFT JOIN whitelist ` +
		`ON categories\.categoryId \= whitelist\.category ` +
		`WHERE categories\.id \= \? AND categories\.categoryId IN \(\?\);`
	mock.
		ExpectPrepare(query1).
		ExpectExec().
		WithArgs(user.Id, budget.Request.Remove.Categories[0].Id).
		WillReturnResult(sqlmock.NewResult(0, 3))
	
	query2 := `DELETE FROM whitelist WHERE whitelist\.id \= \? AND whitelist\.itemId IN \(\?\);`
	mock.
		ExpectPrepare(query2).
		ExpectExec().
		WithArgs(user.Id, budget.Request.Remove.WhiteList[2].Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := models.Delete(app, user.Id, budget)
	test.ModelMethod(t, err, "delete")
	test.MockExpectations(t, mock)
}