package plaid_test

import (
	"encoding/json"
	"net/http"
	"testing"

	m "github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/models"
	p "github.com/elopez00/scale-backend/cmd/api/sdk/plaid"
	"github.com/elopez00/scale-backend/pkg/test"
)

func TestInvalidClient(t *testing.T) {
	app, _ := test.GetMockApp() 
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

		if response.Message != "Failure to load client" {
			t.Errorf("Link token shouldn't have been extracted, instead recieved error: %v", response.Message)
		}
	}
}

func TestLinkTokenRetrieval(t *testing.T) {
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

		if response.Status != 0 {
			t.Errorf("Link token was not extracted successfuly, instead recieved error: %v", response.Message)
		}
	}
}