package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kokoichi206/go-git-stats/api"
	"github.com/kokoichi206/go-git-stats/util"
	"github.com/stretchr/testify/require"
)

func TestListPublicRepositories(t *testing.T) {

	s := httptest.NewServer(nil)
	defer s.Close()

	ts := TestServer{
		server:    s,
		header:    nil,
		apiCalled: 0,
	}

	config := util.Config{
		ApiBaseURL: ts.server.URL,
	}
	// a := api.New(config)
	a := api.ExportNewApi(config)

	testCases := []struct {
		name      string
		userName  string
		setup     func(testServer *httptest.Server)
		assertion func(t *testing.T, err error, repositories []api.Repository)
		tearDown  func()
	}{
		{
			name:     "OK",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusOK, mockRepositories)
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.NoError(t, err)
				t.Log(repositories)
				require.Equal(t, 5, len(repositories))

				expectedRepoNames := []string{"account-book-api", "account-book-ios", "action-URL-watcher", "actions-diff-and-notify", "ai_web_app_flask"}
				for i := 0; i < 5; i++ {
					repository := repositories[i]
					repo := expectedRepoNames[i]
					require.False(t, repository.Private)
					require.Equal(t, repo, repository.Name)
					require.Equal(t, fmt.Sprintf("kokoichi206/%s", repo), repository.FullName)
				}
				t.Log(ts.url)
				require.Equal(t, "/users/kokoichi206/repos", ts.url.Path)

				// Assert header
				passedHeader := ts.header
				t.Log(passedHeader)
				require.NotNil(t, passedHeader)
				hasAccept := false
				for key, values := range passedHeader {
					if key == "Accept" && values[0] == "application/vnd.github+json" {
						hasAccept = true
					}
				}
				require.True(t, hasAccept)

				// Api was called only once
				require.Equal(t, 1, ts.apiCalled)
			},
			tearDown: func() {
			},
		},
		{
			name:     "Error NewRequest",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				t.Log(a.ExportGetConfig())
				a.ExportSetURL("https://test.ser ver.com")
				t.Log(a.ExportGetConfig())
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				t.Log(err)
				require.True(t, strings.Contains(err.Error(), "http.NewRequest"))
				t.Log(repositories)
				require.Nil(t, repositories)

				// Api was NOT called
				require.Equal(t, 0, ts.apiCalled)
			},
			tearDown: func() {
				a.ExportSetURL(ts.server.URL)
			},
		},
		{
			name:     "Error not http scheme",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				t.Log(a.ExportGetConfig())
				a.ExportSetURL("slack://should.return.client.do.err")
				t.Log(a.ExportGetConfig())
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				t.Log(err)
				require.Equal(t, "failed to client.Do after several retries.", err.Error())

				// Api was NOT called
				require.Equal(t, 0, ts.apiCalled)
			},
			tearDown: func() {
				a.ExportSetURL(ts.server.URL)
			},
		},
		{
			name:     "Error Not Found",
			userName: "notFoundUser",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusNotFound, "")
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "client.Do"))
				t.Log(repositories)
				require.Nil(t, repositories)

				t.Log(ts.url)
				require.Equal(t, "/users/notFoundUser/repos", ts.url.Path)

				// Api was called only once
				require.Equal(t, 1, ts.apiCalled)
			},
			tearDown: func() {
			},
		},
		{
			name:     "Error after retry",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusInternalServerError, "")
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "failed to client.Do after several retries."))

				// Api was called 3 times! (and failed...)
				require.Equal(t, 3, ts.apiCalled)
			},
			tearDown: func() {
			},
		},
		{
			name:     "Unmarshal failed with incomplete data",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusOK, mockRepositoriesWithEmpty)
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "json.Unmarshal"))
				t.Log(repositories)
				require.Nil(t, repositories)

				t.Log(ts.url)
				require.Equal(t, "/users/kokoichi206/repos", ts.url.Path)

				// Api was called only once
				require.Equal(t, 1, ts.apiCalled)
			},
			tearDown: func() {
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			tc.setup(ts.server)
			defer tc.tearDown()
			defer ts.init()

			// Act
			repositories, err := a.ListPublicRepositories(tc.userName)

			// Assert
			tc.assertion(t, err, repositories)
		})
	}
}

func TestListRepositoriesForAuthenticatedUser(t *testing.T) {

	s := httptest.NewServer(nil)
	defer s.Close()

	ts := TestServer{
		server:    s,
		header:    nil,
		apiCalled: 0,
	}

	config := util.Config{
		ApiBaseURL: ts.server.URL,
		Token:      "ghq_kokoichi206token",
	}
	a := api.ExportNewApi(config)

	testCases := []struct {
		name      string
		userName  string
		setup     func(testServer *httptest.Server)
		assertion func(t *testing.T, err error, repositories []api.Repository)
		tearDown  func()
	}{
		{
			name:     "OK",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusOK, mockRepositories)
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.NoError(t, err)
				t.Log(repositories)
				require.Equal(t, 5, len(repositories))

				expectedRepoNames := []string{"account-book-api", "account-book-ios", "action-URL-watcher", "actions-diff-and-notify", "ai_web_app_flask"}
				for i := 0; i < 5; i++ {
					repository := repositories[i]
					repo := expectedRepoNames[i]
					require.False(t, repository.Private)
					require.Equal(t, repo, repository.Name)
					require.Equal(t, fmt.Sprintf("kokoichi206/%s", repo), repository.FullName)
				}

				// Assert header
				passedHeader := ts.header
				t.Log(passedHeader)
				require.NotNil(t, passedHeader)
				hasAccept := false
				hasToken := false
				for key, values := range passedHeader {
					if key == "Accept" && values[0] == "application/vnd.github+json" {
						hasAccept = true
					}
					if key == "Authorization" && values[0] == "token ghq_kokoichi206token" {
						hasToken = true
					}
				}
				require.True(t, hasAccept)
				require.True(t, hasToken)

				// Api was called only once
				require.Equal(t, 1, ts.apiCalled)
			},
			tearDown: func() {
			},
		},
		{
			name:     "Error NewRequest",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				a.ExportSetURL("https://test.ser ver.com")
				ts.server.Config.Handler = ts.NewRouter(http.StatusNotFound, "")
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				t.Log(err)
				require.True(t, strings.Contains(err.Error(), "http.NewRequest"))
				t.Log(repositories)
				require.Nil(t, repositories)

				// Api was NOT called
				require.Equal(t, 0, ts.apiCalled)
			},
			tearDown: func() {
				a.ExportSetURL(ts.server.URL)
			},
		},
		{
			name:     "Error not http scheme",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				t.Log(a.ExportGetConfig())
				a.ExportSetURL("slack://should.return.client.do.err")
				t.Log(a.ExportGetConfig())
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				t.Log(err)
				require.Equal(t, "failed to client.Do after several retries.", err.Error())
			},
			tearDown: func() {
				a.ExportSetURL(ts.server.URL)
			},
		},
		{
			name:     "Error Not Found",
			userName: "notFoundUser",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusNotFound, "")
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "client.Do"))
				t.Log(repositories)
				require.Nil(t, repositories)

				// Api was called only once
				require.Equal(t, 1, ts.apiCalled)
			},
			tearDown: func() {
			},
		},
		{
			name:     "Error after retry",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusInternalServerError, "")
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "failed to client.Do after several retries."))

				// Api was called 3 times! (and failed ...)
				require.Equal(t, 3, ts.apiCalled)
			},
			tearDown: func() {
			},
		},
		{
			name:     "Unmarshal failed with incomplete data",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusOK, mockRepositoriesWithEmpty)
			},
			assertion: func(t *testing.T, err error, repositories []api.Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "json.Unmarshal"))
				t.Log(repositories)
				require.Nil(t, repositories)

				// Api was called only once
				require.Equal(t, 1, ts.apiCalled)

			},
			tearDown: func() {
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			tc.setup(ts.server)
			defer tc.tearDown()
			defer ts.init()

			// Act
			repositories, err := a.ListRepositoriesForAuthenticatedUser()

			// Assert
			tc.assertion(t, err, repositories)
		})
	}
}
