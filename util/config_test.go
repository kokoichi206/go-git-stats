package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {

	testCases := []struct {
		name      string
		setup     func()
		assertion func(t *testing.T, config Config)
		tearDown  func()
	}{
		{
			name: "OK",
			setup: func() {
				os.Setenv("GGS_TOKEN", "ghq_foobarTOKEN")
			},
			assertion: func(t *testing.T, config Config) {
				t.Log(config)
				// GitHub REST API
				require.Equal(t, "https://api.github.com", config.ApiBaseURL)
				require.Equal(t, "ghq_foobarTOKEN", config.Token)
			},
			tearDown: func() {
				os.Unsetenv("GGS_TOKEN")
			},
		},
		{
			name:  "OK without token",
			setup: func() {},
			assertion: func(t *testing.T, config Config) {
				t.Log(config)
				// GitHub REST API
				require.Equal(t, "https://api.github.com", config.ApiBaseURL)
				require.Equal(t, "", config.Token)
			},
			tearDown: func() {},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			tc.setup()
			defer tc.tearDown()

			// Act
			config := LoadConfig()

			// Assert
			tc.assertion(t, config)
		})
	}
}
