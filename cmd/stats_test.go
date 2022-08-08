package cmd

import (
	"errors"
	"testing"

	"github.com/kokoichi206/go-git-stats/api"
	"github.com/kokoichi206/go-git-stats/api/mock"
	"github.com/kokoichi206/go-git-stats/util"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestStatsCommand(t *testing.T) {

	config := util.LoadConfig()
	mockApi := mock.New(config)

	c := Cmd{
		Config: config,
		Api:    mockApi,
	}

	app := cli.NewApp()
	app.Commands = c.NewCommands()

	testCases := []struct {
		name      string
		commands  []string
		setup     func()
		assertion func(t *testing.T, err error, api *mock.MockApi)
		tearDown  func()
	}{
		{
			name:     "OK",
			commands: []string{"", "stats", "-name", "kokoichi206/go-git-stats"},
			setup: func() {
				mockApi.ListCodeFreq = []api.CodeFrequency{
					{
						Time:      1659830400,
						Additions: 10_0000,
						Deletions: 999,
					},
				}
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi) {
				require.NoError(t, err)
				require.True(t, api.WeeklyCodeCalled)
				require.Equal(t, "kokoichi206/go-git-stats", api.PassedFullName)
			},
			tearDown: func() {
				mockApi.InitMock()
			},
		},
		{
			name:     "No fullName",
			commands: []string{"", "stats"},
			setup:    func() {},
			assertion: func(t *testing.T, err error, api *mock.MockApi) {
				require.Error(t, err)
				require.False(t, api.WeeklyCodeCalled)
			},
			tearDown: func() {
				mockApi.InitMock()
			},
		},
		{
			name:     "API call error",
			commands: []string{"", "stats", "-name", "kokoichi206/go-git-stats"},
			setup: func() {
				mockApi.Error = errors.New("mock Error: No fullName test")
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi) {
				require.Error(t, err)
				require.True(t, api.WeeklyCodeCalled)
			},
			tearDown: func() {
				mockApi.InitMock()
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
			err := app.Run(tc.commands)

			// Assert
			tc.assertion(t, err, mockApi)
		})
	}
}
