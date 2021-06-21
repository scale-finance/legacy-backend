package middleware_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/sdk"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/test"
)

func TestInvalidAuthentication(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	if res := test.Get("/v0", middleware.Authenticate(sdk.AuthCheck(), app)); res.Code != http.StatusUnauthorized {
		t.Errorf("Wrong http status. Expected %v, got: %v", http.StatusUnauthorized, res.Code)
	} else {
		// Decode response body
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		if response.Message != "Unauthorized User" {
			t.Errorf("This response shouldn't have been possible, expected unauthorized user, got: %v", response.Message)
		}
	}
}

func TestValidAuthentication(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	if res := test.GetWithCookie("/v0", middleware.Authenticate(sdk.AuthCheck(), app), app, "AuthToken"); res.Code != http.StatusOK {
		t.Errorf("Wrong http status. Expected %v, got: %v", http.StatusOK, res.Code)
	} else {
		// encode response body
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		if response.Message != "This app is authenticated" {
			t.Errorf("This response was supposed to work, expected greeting, got: %v", response.Message)
		}
	}
}