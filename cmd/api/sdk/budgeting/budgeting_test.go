package budgeting_test

import (
	"testing"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elopez00/scale-backend/pkg/test"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/sdk/budgeting"
)

var user = models.User {
	Id: "testvalue",
}

var budget = models.Budget{
	Categories: []models.Category {
		{ Name: "shopping", Budget: 200 },
		{ Name: "groceries", Budget: 250 },
		{ Name: "rent", Budget: 800 },
	},

	WhiteList: []models.WhiteListItem {
		{ Category: "shopping", Name: "Calvin Klien" },
		{ Category: "shopping", Name: "Best Buy" },
		{ Category: "shopping", Name: "Amazon" },
		{ Category: "groceries", Name: "Aldi" },
		{ Category: "groceries", Name: "Walmart" },
		{ Category: "rent", Name: "The Rise" },
	},
}

func TestCreateBudget(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	jsonObject, _ := json.Marshal(budget)

	// test categories query
	query1 := `INSERT INTO categories\(id, name, budget\) VALUES \(\?,\?,\?\), \(\?,\?,\?\), \(\?,\?,\?\)`
	mock.ExpectPrepare(query1).
	ExpectExec().
	WillReturnResult(sqlmock.NewResult(0, 0))
	
	// test whitelist query
	query2 := `INSERT INTO whitelist\(id, category, name\) VALUES \(\?,\?,\?\), \(\?,\?,\?\), \(\?,\?,\?\), \(\?,\?,\?\), \(\?,\?,\?\), \(\?,\?,\?\)`
	mock.ExpectPrepare(query2).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))
	
	res := test.PostWithCookie(
		"/v0/createBudget",
		middleware.Authenticate(budgeting.Create(app), app),
		bytes.NewBuffer(jsonObject),
		app,
		"AuthToken",
	)

	if res.Code != http.StatusOK {
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Expected status to be %v, instead we got %v with error message: %v", http.StatusOK, res.Code, response.Message)
		return
	}
}