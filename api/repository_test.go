package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kokoichi206/go-git-stats/util"
	"github.com/stretchr/testify/require"
)

func TestListPublicRepositories(t *testing.T) {

	s := httptest.NewServer(nil)
	defer s.Close()

	ts := TestServer{
		server: s,
		header: nil,
	}

	config := util.Config{
		ApiBaseURL: ts.server.URL,
	}
	api := Api{
		config: config,
	}

	testCases := []struct {
		name      string
		userName  string
		setup     func(testServer *httptest.Server)
		assertion func(t *testing.T, err error, repositories []Repository)
		tearDown  func()
	}{
		{
			name:     "OK",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusOK, mockRepositories)
			},
			assertion: func(t *testing.T, err error, repositories []Repository) {

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
			},
			tearDown: func() {
			},
		},
		{
			name:     "Error NewRequest",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				t.Log(api.config.ApiBaseURL)
				api.config.ApiBaseURL = "https://test.ser ver.com"
				t.Log(api.config.ApiBaseURL)
			},
			assertion: func(t *testing.T, err error, repositories []Repository) {

				require.Error(t, err)
				t.Log(err)
				require.True(t, strings.Contains(err.Error(), "http.NewRequest"))
				t.Log(repositories)
				require.Nil(t, repositories)
			},
			tearDown: func() {
				api.config.ApiBaseURL = ts.server.URL
			},
		},
		{
			name:     "Error Not Found",
			userName: "notFoundUser",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusNotFound, "")
			},
			assertion: func(t *testing.T, err error, repositories []Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "client.Do"))
				t.Log(repositories)
				require.Nil(t, repositories)

				t.Log(ts.url)
				require.Equal(t, "/users/notFoundUser/repos", ts.url.Path)
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
			assertion: func(t *testing.T, err error, repositories []Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "json.Unmarshal"))
				t.Log(repositories)
				require.Nil(t, repositories)

				t.Log(ts.url)
				require.Equal(t, "/users/kokoichi206/repos", ts.url.Path)
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

			// Act
			repositories, err := api.ListPublicRepositories(tc.userName)

			// Assert
			tc.assertion(t, err, repositories)
		})
	}
}

func TestListRepositoriesForAuthenticatedUser(t *testing.T) {

	s := httptest.NewServer(nil)
	defer s.Close()

	ts := TestServer{
		server: s,
		header: nil,
	}

	config := util.Config{
		ApiBaseURL: ts.server.URL,
		Token:      "ghq_kokoichi206token",
	}
	api := Api{
		config: config,
	}

	testCases := []struct {
		name      string
		userName  string
		setup     func(testServer *httptest.Server)
		assertion func(t *testing.T, err error, repositories []Repository)
		tearDown  func()
	}{
		{
			name:     "OK",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusOK, mockRepositories)
			},
			assertion: func(t *testing.T, err error, repositories []Repository) {

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
			},
			tearDown: func() {
			},
		},
		{
			name:     "Error NewRequest",
			userName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				api.config.ApiBaseURL = "https://test.ser ver.com"
				ts.server.Config.Handler = ts.NewRouter(http.StatusNotFound, "")
			},
			assertion: func(t *testing.T, err error, repositories []Repository) {

				require.Error(t, err)
				t.Log(err)
				require.True(t, strings.Contains(err.Error(), "http.NewRequest"))
				t.Log(repositories)
				require.Nil(t, repositories)
			},
			tearDown: func() {
				api.config.ApiBaseURL = ts.server.URL
			},
		},
		{
			name:     "Error Not Found",
			userName: "notFoundUser",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusNotFound, "")
			},
			assertion: func(t *testing.T, err error, repositories []Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "client.Do"))
				t.Log(repositories)
				require.Nil(t, repositories)
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
			assertion: func(t *testing.T, err error, repositories []Repository) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "json.Unmarshal"))
				t.Log(repositories)
				require.Nil(t, repositories)
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

			// Act
			t.Log(api)
			repositories, err := api.ListRepositoriesForAuthenticatedUser()

			// Assert
			tc.assertion(t, err, repositories)
		})
	}
}
