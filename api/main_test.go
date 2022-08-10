package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

type TestServer struct {
	server *httptest.Server
	header http.Header
	url    *url.URL
}

func (ts *TestServer) NewRouter(statusCode int, mockData string) http.Handler {
	fmt.Println("ehhorahora")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(mockData))

		ts.url = r.URL
		// Save passed header
		ts.header = r.Header
	})

	return handler
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
