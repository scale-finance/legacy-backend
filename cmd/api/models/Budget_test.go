package models_test

import (
	"testing"
	"fmt"

	"github.com/elopez00/scale-backend/pkg/test"
	"github.com/elopez00/scale-backend/cmd/api/models"

	"github.com/DATA-DOG/go-sqlmock"
)

var testBudget = models.Budget {
	Categories: []models.Category{
		{Name: "shopping", Budget: 200, WhiteList: []models.WhiteListItem{
			{Category: "shopping", Name: "Calvin Klien"},
			{Category: "shopping", Name: "Best Buy"},
			{Category: "shopping", Name: "Amazon"},
		}},
		{Name: "groceries", Budget: 250, WhiteList: []models.WhiteListItem{
			{Category: "groceries", Name: "Aldi"},
			{Category: "groceries", Name: "Walmart"},
		}},
		{Name: "rent", Budget: 800, WhiteList: []models.WhiteListItem{{Category: "rent", Name: "The Rise"}}},
	},
}

var testBudgetUpdate = models.Budget {
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

// * Test Success

func TestUpdateWhiteList(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := 
		`INSERT INTO whitelist\(id, category, name, itemId\) ` + 
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	
	budget := models.Budget {
		Request: models.UpdateRequest {
			Update: models.UpdateObject {
				WhiteList: testBudgetUpdate.Request.Update.WhiteList,
			},
		},
	}

	whitelist := budget.Request.Update.WhiteList

	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(
			user.Id, whitelist[0].Category, whitelist[0].Name, whitelist[0].Id, 
			user.Id, whitelist[1].Category, whitelist[1].Name, whitelist[1].Id, 
			user.Id, whitelist[2].Category, whitelist[2].Name, whitelist[2].Id, 
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
				Categories: testBudgetUpdate.Request.Update.Categories,
			},
		},
	}

	categories := budget.Request.Update.Categories
	
	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(
			user.Id, categories[0].Name, categories[0].Budget, categories[0].Id,
			user.Id, categories[1].Name, categories[1].Budget, categories[1].Id,	
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := budget.Update(app, user.Id)
	test.ModelMethod(t, err, "insert")
	test.MockExpectations(t, mock)
}

func TestUpdateBudget(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	categories := testBudgetUpdate.Request.Update.Categories
	whitelist := testBudgetUpdate.Request.Update.WhiteList
	
	query2 := 
		`INSERT INTO categories\(id, name, budget, categoryId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId;`
	mock.
		ExpectPrepare(query2).
		ExpectExec().
		WithArgs(
			user.Id, categories[0].Name, categories[0].Budget, categories[0].Id,
			user.Id, categories[1].Name, categories[1].Budget, categories[1].Id,
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
			user.Id, whitelist[0].Category, whitelist[0].Name, whitelist[0].Id,
			user.Id, whitelist[1].Category, whitelist[1].Name, whitelist[1].Id,
			user.Id, whitelist[2].Category, whitelist[2].Name, whitelist[2].Id,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := testBudgetUpdate.Update(app, user.Id)
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
					{ "Shopping", 400, "cid123", []models.WhiteListItem {} },
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
					{ "Shopping", 400, "cid123", []models.WhiteListItem {} },
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

func TestGetBudget(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	categories := testBudget.Categories
	whitelist := []models.WhiteListItem { 
		categories[0].WhiteList[0],
		categories[0].WhiteList[1],
		categories[1].WhiteList[0],
	}

	rows1 := sqlmock.NewRows([]string{"id", "name", "budget", "categoryId"}).
		AddRow(user.Id, categories[0].Name, categories[0].Budget, categories[0].Id).
		AddRow(user.Id, categories[1].Name, categories[1].Budget, categories[1].Id)

	query1 := fmt.Sprintf("SELECT id, name, budget, categoryId FROM categories WHERE categories.id \\= %q", user.Id)
	mock.
		ExpectQuery(query1).
		WillReturnRows(rows1)
	
	rows2 := sqlmock.NewRows([]string{"id", "name", "category", "itemId"}).
		AddRow(user.Id, whitelist[0].Name, whitelist[0].Category, whitelist[0].Id).
		AddRow(user.Id, whitelist[1].Name, whitelist[1].Category, whitelist[1].Id).
		AddRow(user.Id, whitelist[2].Name, whitelist[2].Category, whitelist[2].Id)

	query2 := fmt.Sprintf("SELECT id, name, category, itemId FROM whitelist WHERE whitelist.id \\= %q", user.Id)
	mock.
		ExpectQuery(query2).
		WillReturnRows(rows2)
		
	budget, err := models.GetBudget(app, user.Id)
	test.ModelMethod(t, err, "select")
	test.MockExpectations(t, mock)

	if len(budget.Categories) == 0 {
		t.Fatal("This function was executed successfully, however, it did not return expected budget")
	}

	if budget.Categories[0].Id != testBudget.Categories[0].Id {
		t.Fatal("The function was successfully executed, but the wrong values were returned")
	}
}

// * Test Failure

func TestUpdateFailure(t *testing.T) {
	app, mock := test.GetMockApp()
	app.DB.Client.Close()

	query2 := 
		`INSERT INTO categories\(id, name, budget, categoryId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId;`
	mock.
		ExpectPrepare(query2).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))

	query1 := 
		`INSERT INTO whitelist\(id, category, name, itemId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	mock.
		ExpectPrepare(query1).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))
	
	testBudget.Update(app, user.Id)
	test.MockFailure(t, mock)
}

func TestDeleteFailure(t *testing.T) {
	app, mock := test.GetMockApp()
	app.DB.Client.Close()

	budget := models.Budget {
		Request: models.UpdateRequest {  
			Remove: models.UpdateObject {
				Categories: []models.Category {
					{ "Shopping", 400, "cid123", []models.WhiteListItem {} },
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
	
	models.Delete(app, user.Id, budget)
	test.MockFailure(t, mock)
}

func TestGetBudgetFailure(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	categories := testBudget.Categories
	whitelist := []models.WhiteListItem { 
		categories[0].WhiteList[0],
		categories[0].WhiteList[1],
		categories[1].WhiteList[0],
	}

	rows1 := sqlmock.NewRows([]string{"id", "name", "budget", "categoryId"}).
		AddRow(user.Id, categories[0].Name, categories[0].Budget, categories[0].Id).
		AddRow(user.Id, categories[1].Name, categories[1].Budget, categories[1].Id)

	query1 := fmt.Sprintf("SELECT id, name, budget, categoryId FROM categories WHERE categories.id \\= %q", user.Id)
	mock.
		ExpectQuery(query1).
		WillReturnRows(rows1)
	
	rows2 := sqlmock.NewRows([]string{"id", "name", "category", "itemId"}).
		AddRow(user.Id, whitelist[0].Name, whitelist[0].Category, whitelist[0].Id).
		AddRow(user.Id, whitelist[1].Name, whitelist[1].Category, whitelist[1].Id).
		AddRow(user.Id, whitelist[2].Name, whitelist[2].Category, whitelist[2].Id)

	query2 := fmt.Sprintf("SELECT id, name, category, itemId FROM whitelist WHERE whitelist.id \\= %q", user.Id)
	mock.
		ExpectQuery(query2).
		WillReturnRows(rows2)
		
	_, err := models.GetBudget(app, "sup")
	test.ModelMethodFailure(t, err)
	test.MockFailure(t, mock)
}