package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

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