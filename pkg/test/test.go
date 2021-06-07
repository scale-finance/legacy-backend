package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"testing"
	"time"

	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/cmd/api/sdk/auth"
	"github.com/elopez00/scale-backend/cmd/api/models"
	
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

// This will get a mock application with no live secrets or codes so that the database,
// and general API functions can be tested. It will return a test application and a mock
// database to test queries.
func GetMockApp() (*application.App, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	
	app := application.GetTest(db)

	return app, mock	
}

// This will get a mock application that gets the sandbox keys for plaid. Everything else
// will be returned just as if GetMockApp() were called.
func GetPlaidMockApp() (*application.App, sqlmock.Sqlmock) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		panic(err.Error())
	}
	db, mock, _ := sqlmock.New()

	app := application.GetTest(db)
	return app, mock
}

// This function is used to test post calls with JSON bodies
func Post(endpoint string, handler httprouter.Handle, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", endpoint, body)

	mux := httprouter.New()
	mux.POST(endpoint, handler)

	res := executeRequest(req, mux) 

	return res
}

// This function is used to test get requests without json bodies
func Get(endpoint string, handler httprouter.Handle) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", endpoint, nil)

	mux := httprouter.New()
	mux.GET(endpoint, handler)

	res := executeRequest(req, mux)

	return res
}

// this function is used to test any get request that requires a specific type of cookie. The name
// parameter in this function will be used to specify what cookie the request will search for and
// it will always return a cookie with "testvalue" as its value. Since it is a GET request, this 
// function does not take JSON bodies
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

	res := executeRequest(req, mux)
	return res
}

// This function is used to test any post request that requires a specific type of cookie. The name
// parameter in this function will be used to specify what cookie the request will search for and
// it will always return a cookie with "testvalue" as its value. Since it is a POST request, this
// function will take in a JSON body
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

	res := executeRequest(req, mux)
	return res
}

// Function that will take in the testing object and the mock 
// used for database testing and return a testing error if the 
// expectations were not met for the given mock.
func MockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("There were unfulfilled expectations:", err)
		return
	}
}

// Will take in a testing object, a response recorder, and an expected status code.
// If the code from the response does not match the code from the expected parameter,
// then this function will return a testing error with the response message, and a
// comparison of the response codes.
func Response(t* testing.T, res *httptest.ResponseRecorder, expected int) {
	if res.Code != expected {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Expected %v, got %v, with an error message: %v", expected, res.Code, response.Message)
		return
	}
}

// Given a method, an error, and a testing object, this function 
// will determine if the error is not nil and then return a 
// testing error with a description that will correspond to the 
// described method in the parameter
func ModelMethod(t* testing.T, err error, method string) {
	if err != nil {
		switch(method) {
		case "insert": t.Error("There was an error inserting rows to the database:", err); return
		case "select": t.Error("There was an error getting the rows from the database:", err); return
		case "delete": t.Error("There was an error deleting rows from the database:", err); return
		}
	}
}

// This function will execute any request
func executeRequest(req *http.Request, handler *httprouter.Router) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}