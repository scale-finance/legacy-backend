package models_test

import (
	"testing"

	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/test"

	"github.com/DATA-DOG/go-sqlmock"
)

var testBudget = models.Budget {
	Categories: []models.Category{
		{Name: "shopping", Budget: 200, Color: "red", WhiteList: []models.WhiteListItem{
			{Category: "shopping", Name: "Calvin Klein"},
			{Category: "shopping", Name: "Best Buy"},
			{Category: "shopping", Name: "Amazon"},
		}},
		{Name: "groceries", Budget: 250, Color: "blue", WhiteList: []models.WhiteListItem{
			{Category: "groceries", Name: "Aldi"},
			{Category: "groceries", Name: "Walmart"},
		}},
		{Name: "rent", Budget: 800, Color: "yellow", WhiteList: []models.WhiteListItem{{Category: "rent", Name: "The Rise"}}},
	},
}

var testBudgetUpdate = models.Budget {
	Request: models.UpdateRequest {
		Update: models.UpdateObject {
			Categories: []models.Category { 
				{ Name: "Shopping", Budget: 300, Id: "a;sldfkdj" },
				{ Name: "Fast Food", Budget: 100, Id: "a;sldfkjs" },
			},
			WhiteList: []models.WhiteListItem {
				{ Name: "Polo Store", Category: "Shopping", Id: ";lkj3lk" },
				{ Name: "Five Guys", Category: "Fast Food", Id: ";lkj;fl" },
				{ Name: "Chipotle", Category: "Fast Food", Id: "a;sdl6k" },
			},
		},
	},
}

// * Test Success

func TestUpdateWhiteList(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

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

	app.DB.Mock.
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
	test.MockExpectations(t, app)
}

func TestUpdateCategory(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)
	
	query := 
		`INSERT INTO categories\(id, name, budget, categoryId, color\) ` +
		`VALUES \(\?,\?,\?,\?,\?\), \(\?,\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId, color\=updated\.color;`
	
	budget := models.Budget {
		Request: models.UpdateRequest {
			Update: models.UpdateObject {
				Categories: testBudgetUpdate.Request.Update.Categories,
			},
		},
	}

	categories := budget.Request.Update.Categories
	
	app.DB.Mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(
			user.Id, categories[0].Name, categories[0].Budget, categories[0].Id, categories[0].Color,
			user.Id, categories[1].Name, categories[1].Budget, categories[1].Id, categories[1].Color,	
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := budget.Update(app, user.Id)
	test.ModelMethod(t, err, "insert")
	test.MockExpectations(t, app)
}

func TestUpdateBudget(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	categories := testBudgetUpdate.Request.Update.Categories
	whitelist := testBudgetUpdate.Request.Update.WhiteList
	
	query2 := 
		`INSERT INTO categories\(id, name, budget, categoryId, color\) ` +
		`VALUES \(\?,\?,\?,\?,\?\), \(\?,\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId, color\=updated\.color;`
	app.DB.Mock.
		ExpectPrepare(query2).
		ExpectExec().
		WithArgs(
			user.Id, categories[0].Name, categories[0].Budget, categories[0].Id, categories[0].Color,
			user.Id, categories[1].Name, categories[1].Budget, categories[1].Id, categories[1].Color,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	query1 := 
		`INSERT INTO whitelist\(id, category, name, itemId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	app.DB.Mock.
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
	test.MockExpectations(t, app)
}

func TestDeleteCategoryAndListItems(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	budget := models.Budget {
		Request: models.UpdateRequest {  
			Remove: models.UpdateObject {
				Categories: []models.Category {
					{ "Shopping", 400, "cid123", []models.WhiteListItem {}, "#ff5757" },
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

	app.DB.Mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, budget.Request.Remove.Categories[0].Id).
		WillReturnResult(sqlmock.NewResult(0, 3))
	
	err := models.Delete(app, user.Id, budget)
	test.ModelMethod(t, err, "delete")
	test.MockExpectations(t, app)
}

func TestDeleteWhiteListItem(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

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
	app.DB.Mock.
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
	test.MockExpectations(t, app)
}

func TestDeleteWhiteListAndCategories(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	budget := models.Budget {
		Request: models.UpdateRequest {
			Remove: models.UpdateObject {
				Categories: []models.Category {
					{ "Shopping", 400, "cid123", []models.WhiteListItem {}, "#ff5757" },
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
	app.DB.Mock.
		ExpectPrepare(query1).
		ExpectExec().
		WithArgs(user.Id, budget.Request.Remove.Categories[0].Id).
		WillReturnResult(sqlmock.NewResult(0, 3))
	
	query2 := `DELETE FROM whitelist WHERE whitelist\.id \= \? AND whitelist\.itemId IN \(\?\);`
	app.DB.Mock.
		ExpectPrepare(query2).
		ExpectExec().
		WithArgs(user.Id, budget.Request.Remove.WhiteList[2].Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := models.Delete(app, user.Id, budget)
	test.ModelMethod(t, err, "delete")
	test.MockExpectations(t, app)
}

func TestGetBudget(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	categories := testBudget.Categories
	whitelist := []models.WhiteListItem { 
		categories[0].WhiteList[0],
		categories[0].WhiteList[1],
		categories[1].WhiteList[0],
	}

	rows1 := sqlmock.NewRows([]string{"id", "name", "budget", "categoryId"}).
		AddRow(user.Id, categories[0].Name, categories[0].Budget, categories[0].Id).
		AddRow(user.Id, categories[1].Name, categories[1].Budget, categories[1].Id)

	query1 := `SELECT id, name, budget, categoryId FROM categories WHERE categories.id \= \?`
	app.DB.Mock.
		ExpectQuery(query1).
		WillReturnRows(rows1)
	
	rows2 := sqlmock.NewRows([]string{"id", "name", "category", "itemId"}).
		AddRow(user.Id, whitelist[0].Name, whitelist[0].Category, whitelist[0].Id).
		AddRow(user.Id, whitelist[1].Name, whitelist[1].Category, whitelist[1].Id).
		AddRow(user.Id, whitelist[2].Name, whitelist[2].Category, whitelist[2].Id)

	query2 := `SELECT id, name, category, itemId FROM whitelist WHERE whitelist.id \= \?`
	app.DB.Mock.
		ExpectQuery(query2).
		WillReturnRows(rows2)
		
	budget, err := models.GetBudget(app, user.Id)
	test.ModelMethod(t, err, "select")
	test.MockExpectations(t, app)

	if len(budget.Categories) == 0 {
		t.Fatal("This function was executed successfully, however, it did not return expected budget")
	}

	if budget.Categories[0].Id != testBudget.Categories[0].Id {
		t.Fatal("The function was successfully executed, but the wrong values were returned")
	}
}

// * Test Failure

func TestUpdateFailure(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDBWhenFail(t, app)

	query2 := 
		`INSERT INTO categories\(id, name, budget, categoryId, color\) ` +
		`VALUES \(\?,\?,\?,\?,\?\), \(\?,\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId, color\=updated\.color;`
	app.DB.Mock.
		ExpectPrepare(query2).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))

	query1 := 
		`INSERT INTO whitelist\(id, category, name, itemId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ` +
		`ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	app.DB.Mock.
		ExpectPrepare(query1).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := testBudget.Update(app, user.Id)
	if err != nil {
		return
	}

	test.MockFailure(t, app)
}

func TestDeleteFailure(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDBWhenFail(t, app)

	budget := models.Budget {
		Request: models.UpdateRequest {  
			Remove: models.UpdateObject {
				Categories: []models.Category {
					{ "Shopping", 400, "cid123", []models.WhiteListItem {}, "#ff5757" },
				},
				WhiteList: []models.WhiteListItem {
					{ "cid123", "Calvin Klein", "something" },
					{ "cid123", "Ralph Lauren", "nothing" },
				},
			},
		},
	}

	query := 
		`DELETE categories, whitelist ` +
		`FROM categories LEFT JOIN whitelist ` +
		`ON categories\.categoryId \= whitelist\.category ` +
		`WHERE categories\.id \= \? AND categories\.categoryId IN \(\?\);`

	app.DB.Mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs("not the right id", budget.Request.Remove.Categories[0].Id).
		WillReturnResult(sqlmock.NewResult(0, 3))

	err := models.Delete(app, user.Id, budget)
	if err != nil {
		return
	}

	test.MockFailure(t, app)
}

func TestGetBudgetFailure(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDBWhenFail(t, app)

	categories := testBudget.Categories
	whitelist := []models.WhiteListItem { 
		categories[0].WhiteList[0],
		categories[0].WhiteList[1],
		categories[1].WhiteList[0],
	}

	rows1 := sqlmock.NewRows([]string{"id", "name", "budget", "categoryId"}).
		AddRow(user.Id, categories[0].Name, categories[0].Budget, categories[0].Id).
		AddRow(user.Id, categories[1].Name, categories[1].Budget, categories[1].Id)

	query1 := `SELECT id, name, budget, categoryId FROM categories WHERE categories.id \= \?`
	app.DB.Mock.
		ExpectQuery(query1).
		WillReturnRows(rows1)
	
	rows2 := sqlmock.NewRows([]string{"id", "name", "category", "itemId"}).
		AddRow(user.Id, whitelist[0].Name, whitelist[0].Category, whitelist[0].Id).
		AddRow(user.Id, whitelist[1].Name, whitelist[1].Category, whitelist[1].Id).
		AddRow(user.Id, whitelist[2].Name, whitelist[2].Category, whitelist[2].Id)

	query2 := `SELECT id, name, category, itemId FROM whitelist WHERE whitelist.id \= \something`
	app.DB.Mock.
		ExpectQuery(query2).
		WillReturnRows(rows2)
		
	_, err := models.GetBudget(app, "sup")
	test.ModelMethodFailure(t, err)
	test.MockFailure(t, app)
}