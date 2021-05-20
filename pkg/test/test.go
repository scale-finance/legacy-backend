package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elopez00/scale-backend/cmd/api/models"
	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/pkg/cookie"
	"github.com/julienschmidt/httprouter"
)

func GetMockApp() (*application.App, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	app := application.GetTest(db)
	return app, mock
}

func Post(endpoint string, handler httprouter.Handle, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", endpoint, body)

	mux := httprouter.New()
	mux.POST(endpoint, handler)

	res := ExecuteRequest(req, mux) 

	return res
}

func Get(endpoint string, handler httprouter.Handle, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", endpoint, body)

	mux := httprouter.New()
	mux.GET(endpoint, handler)

	res := ExecuteRequest(req, mux)

	return res
}

func GetWithCookie(endpoint string, handler httprouter.Handle, body io.Reader, app *application.App) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", endpoint, body)
	token, _ := cookie.GenerateJWT(app, "testvalue")
	req.AddCookie(&http.Cookie {
		Name: "testcookie",
		Value: token,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	})

	mux := httprouter.New()
	mux.GET(endpoint, handler)

	res := ExecuteRequest(req, mux)

	return res
}

func ExecuteRequest(req *http.Request, handler *httprouter.Router) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}

func CheckResponse(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func Handler(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.Response {
			Status: 0,
			Type: "Hello",
			Message: "Hello World!",
		})
	}
}