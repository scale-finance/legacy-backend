package plaid_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	m "github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/models"
	p "github.com/elopez00/scale-backend/cmd/api/sdk/plaid"
	"github.com/elopez00/scale-backend/pkg/test"
)

var token = models.Token{
	Value: "access-sandbox-3b6a6577-4c02-4fc3-a213-b8adf828c38f",
	Id:    "nothin",
	Institution:  "institution",
}

var publicToken = models.Token{
	Value: "public-sandbox-4d532c06-b9b5-4a18-906a-df480f320cc9",
}

var user = models.User{
	Id: "testvalue",
}

// * Test Functions will invalid Plaid clients *

func TestLinkTokenInvalidClient(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	res := test.GetWithCookie(
		"/v0/getLinkToken",
		m.Authenticate(p.GetPlaidToken(app), app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

func TestExchangePublicTokenInvalidClient(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	body, _ := json.Marshal(publicToken)

	res := test.PostWithCookie(
		"/v0/exchangePublicToken",
		m.Authenticate(p.ExchangePublicToken(app), app),
		bytes.NewBuffer(body),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

func TestGetTransactionsInvalidClient(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	body, _ := json.Marshal(publicToken)

	res := test.PostWithCookie(
		"/v0/exchangePublicToken",
		m.Authenticate(p.GetTransactions(app), app),
		bytes.NewBuffer(body),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

func TestGetBalancesInvalidClient(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	res := test.GetWithCookie(
		"/v0/getBalances",
		m.Authenticate(p.GetBalance(app), app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

// * Test calls with valid Plaid Clients *

func TestGetLinkToken(t *testing.T) {
	app, _ := test.GetPlaidMockApp()
	defer app.DB.Client.Close()

	res := test.GetWithCookie(
		"/v0/getLinkToken",
		m.Authenticate(p.GetPlaidToken(app), app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusOK)
}

func TestGetTransactions(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"id", "token", "itemID", "institution"}).
		AddRow(user.Id, token.Value, token.Id, token.Institution)

	query := `SELECT id, token, itemID, institution FROM plaidtokens WHERE id\="testvalue"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	res := test.GetWithCookie(
		"/v0/getTransactions",
		m.Authenticate(p.GetTransactions(app), app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusOK)
}

func TestGetBalances(t *testing.T) {
	app, mock := test.GetPlaidMockApp()
	defer app.DB.Client.Close()

	// mock the database retrieval
	rows := sqlmock.NewRows([]string{"id", "token", "itemID", "institution"}).
		AddRow(user.Id, token.Value, token.Id, token.Institution)

	query := `SELECT id, token, itemID, institution FROM plaidtokens WHERE id\="testvalue"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	res := test.GetWithCookie(
		"/v0/getBalance",
		m.Authenticate(p.GetBalance(app), app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusOK)
	test.MockExpectations(t, mock)
}

// * Testing error messages

func TestExchangePublicTokenInvalidToken(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	body, _ := json.Marshal(publicToken)

	res := test.PostWithCookie(
		"/v0/exchangePublicToken",
		p.ExchangePublicToken(app),
		bytes.NewBuffer(body),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

func TestGetBalancesInvalidToken(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	res := test.GetWithCookie(
		"/v0/getBalances",
		p.GetBalance(app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}
