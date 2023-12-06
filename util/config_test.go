package util_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kokoichi206/go-git-stats/util"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {

	testCases := []struct {
		name      string
		setup     func()
		assertion func(t *testing.T, config util.Config, err error)
		tearDown  func()
	}{
		{
			name: "OK",
			setup: func() {
				os.Setenv("GGS_TOKEN", "ghp_q5k0u9YD8JPHVIUckyx4dKDyvGavdJWHR44D")
			},
			assertion: func(t *testing.T, config util.Config, err error) {
				t.Log(config)
				// GitHub REST API
				require.Equal(t, "https://api.github.com", config.ApiBaseURL)
				require.Equal(t, "ghp_q5k0u9YD8JPHVIUckyx4dKDyvGavdJWHR44D", config.Token)
				require.NoError(t, err)
			},
			tearDown: func() {
				os.Unsetenv("GGS_TOKEN")
			},
		},
		{
			name:  "OK without token",
			setup: func() {},
			assertion: func(t *testing.T, config util.Config, err error) {
				t.Log(config)
				// GitHub REST API
				require.Equal(t, "https://api.github.com", config.ApiBaseURL)
				require.Equal(t, "", config.Token)
				require.NoError(t, err)
			},
			tearDown: func() {},
		},
		{
			name: "Token format error",
			setup: func() {
				os.Setenv("GGS_TOKEN", "invalid_token")
			},
			assertion: func(t *testing.T, config util.Config, err error) {
				t.Log(config)
				require.Equal(t, "", config.ApiBaseURL)
				require.Equal(t, "", config.Token)
				require.Error(t, err)
				require.Equal(t, err.Error(), fmt.Sprintf("Your token: '%s' is invalid format.\nPlease check your environment variable [GGS_TOKEN].", "invalid_token"))
			},
			tearDown: func() {
				os.Unsetenv("GGS_TOKEN")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			tc.setup()
			defer tc.tearDown()

			// Act
			config, err := util.LoadConfig()

			// Assert
			tc.assertion(t, config, err)
		})
	}
}

func TestIsValidFormat(t *testing.T) {

	testCases := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			name:     "OK",
			token:    "ghp_6nKaaa1xmIdPQ0jU9fZwOMx97qjI6343PzoT",
			expected: true,
		},
		{
			name:     "NG wrong prefix",
			token:    "ghr_6nKaaa1xmIdPQ0jU9fZwOMx97qjI6343PzoT",
			expected: false,
		},
		{
			name:     "NG short",
			token:    "ghp_6nKaaa",
			expected: false,
		},
		{
			name:     "NG wrong prefix",
			token:    "ghr_6nKaaa1xmIdPQ0jU9fZwOMx97qjI6343PzoTKp43jfafeapjdpl",
			expected: false,
		},
		{
			name:     "NG wrong character",
			token:    "ghr_wrong;character;jU9fZwOMxqjI6343PzoT",
			expected: false,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			result := util.IsValidFormat(tc.token)

			// Assert
			require.Equal(t, tc.expected, result)
		})
	}
}
