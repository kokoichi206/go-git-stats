package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kokoichi206/go-git-stats/util"
	"github.com/stretchr/testify/require"
)

func TestWeeklyCommitActivity(t *testing.T) {

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
		fullName  string
		setup     func(testServer *httptest.Server)
		assertion func(t *testing.T, err error, frequencies []CodeFrequency)
		tearDown  func()
	}{
		{
			name:     "OK",
			fullName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusOK, mockCodeFrequencies)
			},
			assertion: func(t *testing.T, err error, frequencies []CodeFrequency) {

				require.NoError(t, err)
				t.Log(frequencies)
				require.Equal(t, 3, len(frequencies))

				// CAUTION: reversed order
				expectedTimes := []int{1627171200, 1626566400, 1625961600}
				expectedAdditions := []int{3375, 23550, 4381719}
				expectedDeletions := []int{-813, -208, -9488}
				for i := 0; i < 3; i++ {
					frequency := frequencies[i]
					require.Equal(t, frequency.Time, expectedTimes[i])
					require.Equal(t, frequency.Additions, expectedAdditions[i])
					require.Equal(t, frequency.Deletions, expectedDeletions[i])
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
			fullName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				api.config.ApiBaseURL = "https://test.ser ver.com"
				ts.server.Config.Handler = ts.NewRouter(http.StatusNotFound, "")
			},
			assertion: func(t *testing.T, err error, frequencies []CodeFrequency) {

				require.Error(t, err)
				t.Log(err)
				require.True(t, strings.Contains(err.Error(), "http.NewRequest"))
				t.Log(frequencies)
				require.Nil(t, frequencies)
			},
			tearDown: func() {
				api.config.ApiBaseURL = ts.server.URL
			},
		},
		{
			name:     "Error Not Found",
			fullName: "notFoundUser",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusNotFound, "")
			},
			assertion: func(t *testing.T, err error, frequencies []CodeFrequency) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "client.Do"))
				t.Log(frequencies)
				require.Nil(t, frequencies)
			},
			tearDown: func() {
			},
		},
		{
			name:     "Unmarshal failed with incomplete data",
			fullName: "kokoichi206",
			setup: func(testServer *httptest.Server) {
				ts.server.Config.Handler = ts.NewRouter(http.StatusOK, codeFrequenciesUnmarshalError)
			},
			assertion: func(t *testing.T, err error, frequencies []CodeFrequency) {

				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), "json.Unmarshal"))
				t.Log(frequencies)
				require.Nil(t, frequencies)
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
			frequencies, err := api.WeeklyCommitActivity(tc.fullName)

			// Assert
			tc.assertion(t, err, frequencies)
		})
	}
}
