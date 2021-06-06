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
}

func TestTokenAdd(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	query := `INSERT INTO plaidtokens\(id, token, itemID\) VALUES\(\?,\?,\?\)`
	mock.
		ExpectPrepare(query).
		ExpectExec().
		WithArgs(user.Id, token.Value, token.Id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := token.Add(app, user.Id); err != nil {
		t.Error("There was an error adding the token to the database: ", err)
		return
	}

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

