package sdk_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	m "github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/sdk"
	"github.com/elopez00/scale-backend/pkg/test"
)

// * Test Functions will invalid Plaid clients *

func TestLinkTokenInvalidClient(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	res := test.GetWithCookie(
		"/v0/getLinkToken",
		m.Authenticate(sdk.GetPlaidToken(app), app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

func TestExchangePublicTokenInvalidClient(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	body, _ := json.Marshal(publicToken)

	res := test.PostWithCookie(
		"/v0/exchangePublicToken",
		m.Authenticate(sdk.ExchangePublicToken(app), app),
		bytes.NewBuffer(body),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

func TestGetTransactionsInvalidClient(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	body, _ := json.Marshal(publicToken)

	res := test.PostWithCookie(
		"/v0/exchangePublicToken",
		m.Authenticate(sdk.GetTransactions(app), app),
		bytes.NewBuffer(body),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

func TestGetBalancesInvalidClient(t *testing.T) {
	app := test.GetMockApp()
	defer test.CloseDB(t, app)

	res := test.GetWithCookie(
		"/v0/getBalances",
		m.Authenticate(sdk.GetBalance(app), app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

// * Test calls with valid Plaid Clients *

func TestGetLinkToken(t *testing.T) {
	app := test.GetPlaidMockApp()
	defer test.CloseDB(t, app)

	res := test.GetWithCookie(
		"/v0/getLinkToken",
		m.Authenticate(sdk.GetPlaidToken(app), app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusOK)
}

//func TestGetTransactions(t *testing.T) {
//	app := test.GetPlaidMockApp()
//	defer test.CloseDB(t, app)
//
//	rows := sqlmock.NewRows([]string{"id", "token", "itemID", "institution"}).
//		AddRow(user.Id, token.Value, token.Id, token.Institution)
//
//	query := `SELECT id, token, itemID, institution FROM plaidtokens WHERE id\="testvalue"`
//	app.DB.Mock.ExpectQuery(query).WillReturnRows(rows)
//
//	res := test.GetWithCookie(
//		"/v0/getTransactions",
//		m.Authenticate(sdk.GetTransactions(app), app),
//		app,
//		"AuthToken",
//	)
//
//	test.Response(t, res, http.StatusOK)
//}

//func TestGetBalances(t *testing.T) {
//	app := test.GetPlaidMockApp()
//	defer test.CloseDB(t, app)
//
//	// mock the database retrieval
//	rows := sqlmock.NewRows([]string{"id", "token", "itemID", "institution"}).
//		AddRow(user.Id, token.Value, token.Id, token.Institution)
//
//	query := `SELECT id, token, itemID, institution FROM plaidtokens WHERE id\="testvalue"`
//	app.DB.Mock.ExpectQuery(query).WillReturnRows(rows)
//
//	res := test.GetWithCookie(
//		"/v0/getBalance",
//		m.Authenticate(sdk.GetBalance(app), app),
//		app,
//		"AuthToken",
//	)
//
//	test.Response(t, res, http.StatusOK)
//	test.MockExpectations(t, app)
//}

// * Testing error messages

func TestExchangePublicTokenInvalidToken(t *testing.T) {
	app := test.GetPlaidMockApp()
	defer test.CloseDB(t, app)

	body, _ := json.Marshal(publicToken)

	res := test.PostWithCookie(
		"/v0/exchangePublicToken",
		sdk.ExchangePublicToken(app),
		bytes.NewBuffer(body),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}

func TestGetBalancesInvalidToken(t *testing.T) {
	app := test.GetPlaidMockApp()
	defer test.CloseDB(t, app)

	res := test.GetWithCookie(
		"/v0/getBalances",
		sdk.GetBalance(app),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusBadGateway)
}
