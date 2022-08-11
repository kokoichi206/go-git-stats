package api_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

type TestServer struct {
	server    *httptest.Server
	header    http.Header
	url       *url.URL
	apiCalled int
}

func (ts *TestServer) NewRouter(statusCode int, mockData string) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(mockData))

		ts.url = r.URL
		// Save passed header
		ts.header = r.Header

		// Check how many times is the API called
		ts.apiCalled += 1
	})

	return handler
}

func (ts *TestServer) init() {
	ts.header = nil
	ts.url = nil
	ts.apiCalled = 0
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
