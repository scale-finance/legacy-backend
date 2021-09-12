package models_test

import (
	"testing"

	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/test"

	"github.com/DATA-DOG/go-sqlmock"
)

var token = models.Token{
	Value: "randomaccess",
	Id:    "randomid",
	Institution:  "Bank of Bank",
}

func TestTokenAdd(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	query := `INSERT INTO plaidtokens\(id, token, itemID, institution\) VALUES\(\?,\?,\?,\?\)`
	app.DB.Mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, token.Value, token.Id, token.Institution).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := token.Add(app, user.Id)
	test.ModelMethod(t, err, "insert")
	test.MockExpectations(t, app)
}

func TestGetTokens(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	rows := sqlmock.NewRows([]string{"id", "token", "itemID", "institution"}).
		AddRow(user.Id, "token1", "id1", "institution1").
		AddRow(user.Id, "token2", "id2", "institution2")

	query := `SELECT id, token, itemID, institution FROM plaidtokens WHERE id \= ?`
	app.DB.Mock.ExpectQuery(query).WillReturnRows(rows)

	tokens, err := models.GetTokens(app, user.Id)
	test.ModelMethod(t, err, "select")
	test.MockExpectations(t, app)

	if len(tokens) == 0 {
		t.Error("The function did not return the tokens")
		return
	}
}

func TestGetToken(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	row := sqlmock.NewRows([]string{"id", "token", "itemID", "institution"}).
		AddRow(user.Id, "token1", "id1", "institution1")

	query := `SELECT id, token, itemID, institution FROM plaidtokens WHERE id \= \? AND itemID \= \?`
	app.DB.Mock.ExpectQuery(query).WillReturnRows(row)

	tempToken := token
	err := tempToken.Get(app, user.Id)

	test.ModelMethod(t, err, "select")
	test.MockExpectations(t, app)

	if len(tempToken.Value) == 0 {
		t.Error("The function did not return the specified token")
		return
	}
}
