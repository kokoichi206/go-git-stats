package cmd

import (
	"testing"

	"github.com/kokoichi206/go-git-stats/api"
	"github.com/kokoichi206/go-git-stats/api/mock"
	"github.com/kokoichi206/go-git-stats/util"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestCommand(t *testing.T) {

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
			name:     "OK with access token",
			commands: []string{"", "repo"},
			setup: func() {
				c.Config.Token = "TokenString"
				mockApi.ListRepos = []api.Repository{
					{
						ID:       489517307,
						Private:  false,
						Name:     "account-book-api",
						FullName: "kokoichi206/account-book-api",
					},
				}
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi) {
				require.NoError(t, err)
				require.False(t, api.PublicCalled)
				require.True(t, api.AuthenticatedCalled)
			},
			tearDown: func() {
				c.Config.Token = ""
				mockApi.InitMock()
			},
		},
		{
			name:     "OK with user name",
			commands: []string{"", "repo", "-name", "kokoichi206"},
			setup: func() {
				c.Config.Token = ""
				mockApi.ListRepos = []api.Repository{
					{
						ID:       489517307,
						Private:  false,
						Name:     "account-book-api",
						FullName: "kokoichi206/account-book-api",
					},
				}
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi) {
				require.NoError(t, err)
				require.True(t, api.PublicCalled)
				require.False(t, api.AuthenticatedCalled)
			},
			tearDown: func() {
				c.Config.Token = ""
				mockApi.InitMock()
			},
		},
		{
			name:     "abbr of subcommand",
			commands: []string{"", "r", "-name", "kokoichi206"},
			setup: func() {
				c.Config.Token = ""
				mockApi.ListRepos = []api.Repository{
					{
						ID:       489517307,
						Private:  false,
						Name:     "account-book-api",
						FullName: "kokoichi206/account-book-api",
					},
				}
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi) {
				require.NoError(t, err)
				require.True(t, api.PublicCalled)
				require.False(t, api.AuthenticatedCalled)
			},
			tearDown: func() {
				c.Config.Token = ""
				mockApi.InitMock()
			},
		},
		{
			name:     "abbr of name flag",
			commands: []string{"", "repo", "-n", "kokoichi206"},
			setup: func() {
				c.Config.Token = ""
				mockApi.ListRepos = []api.Repository{
					{
						ID:       489517307,
						Private:  false,
						Name:     "account-book-api",
						FullName: "kokoichi206/account-book-api",
					},
				}
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi) {
				require.NoError(t, err)
				require.True(t, api.PublicCalled)
				require.False(t, api.AuthenticatedCalled)
			},
			tearDown: func() {
				c.Config.Token = ""
				mockApi.InitMock()
			},
		},
		{
			name:     "Token or userName is not given",
			commands: []string{"", "repo"},
			setup: func() {
				c.Config.Token = ""
				mockApi.ListRepos = []api.Repository{
					{
						ID:       489517307,
						Private:  false,
						Name:     "account-book-api",
						FullName: "kokoichi206/account-book-api",
					},
				}
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi) {
				require.Error(t, err)
				require.Equal(t, "Token or userName is not given.", err.Error())
				require.False(t, api.PublicCalled)
				require.False(t, api.AuthenticatedCalled)
			},
			tearDown: func() {
				c.Config.Token = ""
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
