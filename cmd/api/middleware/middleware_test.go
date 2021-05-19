package middleware_test

import (
	"net/http"
	"testing"
	"encoding/json"

	"github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/test"
)

func TestUnvalidAuthentication(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	if res := test.Get("/", middleware.Authenticate(test.Handler(app), app), nil); res.Code != http.StatusUnauthorized {
		t.Errorf("Wrong http status. Expected %v, got: %v", http.StatusUnauthorized, res.Code)
	} else {
		// encode response body
		var response models.Response
		json.NewEncoder(res).Encode(&response)

		if response.Message != "Unauthorized User" {
			t.Errorf("This response shouldn't have been possible, expected unauthorized user, got: %v", response.Message)
		}
	}
}

func TestValidAuthentication(t *testing.T) {
	app, _ := test.GetMockApp()
	defer app.DB.Client.Close()

	if res := test.Get("/", middleware.Authenticate(test.Handler(app), app), nil); res.Code != http.StatusOK {
		t.Errorf("Wrong http status. Expected %v, got: %v", http.StatusOK, res.Code)
	} else {
		// encode response body
		var response models.Response
		json.NewEncoder(res).Encode(&response)

		if response.Message != "Hello World!" {
			t.Errorf("This response shouldn't have been possible, expected greeting, got: %v", response.Message)
		}
	}
}