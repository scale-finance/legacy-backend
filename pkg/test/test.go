package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/cmd/api/sdk/auth"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

// TODO document this package

func GetMockApp() (*application.App, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	
	app := application.GetTest(db)

	return app, mock	
}

func GetPlaidMockApp() (*application.App, sqlmock.Sqlmock) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		panic(err.Error())
	}
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

func GetWithCookie(endpoint string, handler httprouter.Handle, app *application.App, name string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", endpoint, nil)
	token, _ := auth.GenerateJWT(app, "testvalue")
	req.AddCookie(&http.Cookie {
		Name: name,
		Value: token,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	})

	mux := httprouter.New()
	mux.GET(endpoint, handler)

	res := ExecuteRequest(req, mux)
	return res
}

func PostWithCookie(endpoint string, handler httprouter.Handle, body io.Reader, app *application.App, name string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", endpoint, body)
	token, _ := auth.GenerateJWT(app, "testvalue")
	req.AddCookie(&http.Cookie {
		Name: name,
		Value: token,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	})

	mux := httprouter.New()
	mux.POST(endpoint, handler)

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