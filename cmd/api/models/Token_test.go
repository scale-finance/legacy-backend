package models_test

import (
	"testing"

	"github.com/elopez00/scale-backend/pkg/test"
	"github.com/elopez00/scale-backend/cmd/api/models"

	"github.com/DATA-DOG/go-sqlmock"
)


var token = models.Token{
	Value: "randomaccess",
	Id:    "randomid",
	Name:  "Bank of Bank",
}

func TestTokenAdd(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `INSERT INTO plaidtokens\(id, token, itemID, institution\) VALUES\(\?,\?,\?,\?\)`
	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, token.Value, token.Id, token.Name).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := token.Add(app, user.Id)
	test.ModelMethod(t, err, "insert")
	test.MockExpectations(t, mock)
}

func TestGetTokens(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"id", "token", "itemID", "institution"}).
		AddRow(user.Id, "token1", "id1", "institution1").
		AddRow(user.Id, "token2", "id2", "institution2")

	query := `SELECT id, token, itemID, institution FROM plaidtokens WHERE id\="goingdowntosouthpark"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	tokens, err := models.GetTokens(app, user.Id)
	test.ModelMethod(t, err, "select")
	test.MockExpectations(t, mock)

	if len(tokens) == 0 {
		t.Error("The function did not return the tokens")
		return
	}
}

