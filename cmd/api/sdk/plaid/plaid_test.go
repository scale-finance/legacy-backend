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

	if res := test.GetWithCookie(
		"/v0/getLinkToken",
		m.Authenticate(p.GetPlaidToken(app), app),
		nil,
		app,
		"AuthToken",
	); res.Code != http.StatusBadGateway {
		t.Errorf("Failed get. Expected %v, instead got %v", http.StatusOK, res.Code)
	} else {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		if response.Message != "Failure to load client" {
			t.Errorf("Link token shouldn't have been extracted, instead recieved error: %v", response.Message)
		}
	}
}

func TestExchangePublicTokenInvalidClient(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	body, _ := json.Marshal(publicToken)

	res := test.GetWithCookie(
		"/v0/exchangePublicToken",
		m.Authenticate(p.ExchangePublicToken(app), app),
		bytes.NewBuffer(body),
		app,
		"AuthToken",
	)

	if res.Code != http.StatusBadGateway {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Expected %v, got %v, with an error message: %v", http.StatusBadGateway, res.Code, response.Message)
	}
}

func TestGetTransactionsInvalidClient(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	body, _ := json.Marshal(publicToken)

	res := test.GetWithCookie(
		"/v0/exchangePublicToken",
		m.Authenticate(p.GetTransactions(app), app),
		bytes.NewBuffer(body),
		app,
		"AuthToken",
	)

	if res.Code != http.StatusBadGateway {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Expected %v, got %v, with an error message: %v", http.StatusBadGateway, res.Code, response.Message)
	}
}

func TestGetBalancesInvalidClient(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	res := test.GetWithCookie(
		"/v0/getBalances",
		m.Authenticate(p.GetBalance(app), app),
		nil,
		app,
		"AuthToken",
	)

	if res.Code != http.StatusBadGateway {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Expected %v, got %v, with an error message: %v", http.StatusBadGateway, res.Code, response.Message)
	}
}

// * Test calls with valid Plaid Clients *

func TestGetLinkToken(t *testing.T) {
	app, _ := test.GetPlaidMockApp()
	defer app.DB.Client.Close()

	if res := test.GetWithCookie(
		"/v0/getLinkToken",
		m.Authenticate(p.GetPlaidToken(app), app),
		nil,
		app,
		"AuthToken",
	); res.Code != http.StatusOK {
		t.Errorf("Failed get. Expected %v, instead got %v", http.StatusOK, res.Code)
	} else {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		if response.Message != "Successfully recieved link token from plaid" {
			t.Errorf("Link token was not extracted successfuly, instead recieved error: %v", response.Message)
		}
	}
}

func TestGetTransactions(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"id", "token", "itemID"}).
		AddRow(user.Id, token.Value, token.Id)

	query := `SELECT id, token, itemID FROM plaidtokens WHERE id\="testvalue"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	res := test.GetWithCookie(
		"/v0/getTransactions",
		m.Authenticate(p.GetTransactions(app), app),
		nil,
		app,
		"AuthToken",
	)

	if res.Code != http.StatusOK {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("This call returned the wrong http status. Expected %v, got %v", http.StatusOK, res.Code)
		t.Error("The call did not return the intended result, instead", response.Message)
	}
}

func TestGetBalances(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	// mock the database retrieval
	rows := sqlmock.NewRows([]string{"id", "token", "itemID"}).
		AddRow(user.Id, token.Value, token.Id)

	query := `SELECT id, token, itemID FROM plaidtokens WHERE id\="testvalue"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	res := test.GetWithCookie(
		"/v0/getBalance",
		m.Authenticate(p.GetBalance(app), app),
		nil,
		app,
		"AuthToken",
	)

	var response models.Response
	json.NewDecoder(res.Body).Decode(&response)

	if res.Code != http.StatusOK {
		t.Errorf("There was an error getting account balances, expected 200, got: %v, with error message %v", res.Code, response)
		return
	}

	if response.Result == nil {
		t.Error("The call was successful, but the function did not return a valid response")
		return
	}
}

// * Testing error messages

func TestExchangePublicTokenInvalidToken(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	body, _ := json.Marshal(publicToken)

	res := test.GetWithCookie(
		"/v0/exchangePublicToken",
		p.ExchangePublicToken(app),
		bytes.NewBuffer(body),
		app,
		"AuthToken",
	)

	if res.Code != http.StatusBadGateway {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Expected %v, got %v, with an error message: %v", http.StatusBadGateway, res.Code, response.Message)
	}
}

func TestGetBalancesInvalidToken(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	res := test.GetWithCookie(
		"/v0/getBalances",
		p.GetBalance(app),
		nil,
		app,
		"AuthToken",
	)

	if res.Code != http.StatusBadGateway {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Expected %v, got %v, with an error message: %v", http.StatusBadGateway, res.Code, response.Message)
	}
}