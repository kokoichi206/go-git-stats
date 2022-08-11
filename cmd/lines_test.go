package cmd

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/kokoichi206/go-git-stats/api"
	"github.com/kokoichi206/go-git-stats/api/mock"
	"github.com/kokoichi206/go-git-stats/util"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestLinesCommand(t *testing.T) {

	config, _ := util.LoadConfig()
	mockApi := mock.New(config)

	c := Cmd{
		Config: config,
		Api:    mockApi,
		Wait:   &sync.WaitGroup{},
		Mutex:  &sync.Mutex{},
	}

	app := cli.NewApp()
	app.Commands = c.NewCommands()

	testCases := []struct {
		name      string
		commands  []string
		setup     func()
		assertion func(t *testing.T, err error, api *mock.MockApi, output string)
		tearDown  func()
	}{
		{
			name:     "OK",
			commands: []string{"", "lines", "-n", "kokoichi206"},
			setup: func() {
				c.total = 0
				c.Config.Token = ""
				mockApi.ListRepos = []api.Repository{
					{
						ID:       489517307,
						Private:  false,
						Name:     "account-book-api",
						FullName: "kokoichi206/account-book-api",
					},
					{
						ID:       429817377,
						Private:  false,
						Name:     "utils",
						FullName: "kokoichi206/utils",
					},
				}
				mockApi.ListCodeFreq = append(mockApi.ListCodeFreq, []api.CodeFrequency{
					{
						Time:      1659830400,
						Additions: 9_5000,
						Deletions: 5000,
					},
					{
						Time:      1659225600,
						Additions: 500,
						Deletions: 1,
					},
				})
				mockApi.ListCodeFreq = append(mockApi.ListCodeFreq, []api.CodeFrequency{
					{
						Time:      1659830400,
						Additions: 200,
						Deletions: 800,
					},
				})
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi, output string) {
				require.NoError(t, err)
				require.True(t, api.PublicCalled)
				require.False(t, api.AuthenticatedCalled)
				require.True(t, api.WeeklyCodeCalled)
				require.Equal(t, 101501, c.total)

				// Sum of lines of codes
				t.Log(output)
				require.True(t, strings.Contains(output, "101501"))
			},
			tearDown: func() {
				mockApi.InitMock()
			},
		},
		{
			name:     "OK with token",
			commands: []string{"", "lines"},
			setup: func() {
				c.total = 0
				c.Config.Token = "ghq_foobartoken"
				mockApi.ListCodeFreq = nil
				mockApi.ListRepos = []api.Repository{
					{
						ID:       489517307,
						Private:  false,
						Name:     "account-book-api",
						FullName: "kokoichi206/account-book-api",
					},
					{
						ID:       429817377,
						Private:  false,
						Name:     "utils",
						FullName: "kokoichi206/utils",
					},
				}
				mockApi.ListCodeFreq = append(mockApi.ListCodeFreq, []api.CodeFrequency{
					{
						Time:      1659830400,
						Additions: 9_5000,
						Deletions: 5000,
					},
				})
				mockApi.ListCodeFreq = append(mockApi.ListCodeFreq, []api.CodeFrequency{
					{
						Time:      1659830400,
						Additions: 2000_0000,
						Deletions: 0,
					},
				})
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi, output string) {
				require.NoError(t, err)
				require.False(t, api.PublicCalled)
				require.True(t, api.AuthenticatedCalled)
				require.True(t, api.WeeklyCodeCalled)
				require.Equal(t, 20100000, c.total)

				// Sum of lines of codes
				t.Log(output)
				require.True(t, strings.Contains(output, "20100000"))
			},
			tearDown: func() {
				mockApi.InitMock()
			},
		},
		{
			name:     "Without token and username",
			commands: []string{"", "lines"},
			setup: func() {
				c.total = 0
				c.Config.Token = ""
			},
			assertion: func(t *testing.T, err error, api *mock.MockApi, output string) {
				require.NoError(t, err)
				require.False(t, api.PublicCalled)
				require.False(t, api.AuthenticatedCalled)
				require.False(t, api.WeeklyCodeCalled)
				require.Equal(t, 0, c.total)

				// Sum of lines of codes
				t.Log(output)
				require.True(t, strings.Contains(output, "0"))
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

			// Prepare for standard output testing
			stdOut := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Act
			err := app.Run(tc.commands)

			_ = w.Close()
			result, _ := io.ReadAll(r)
			output := string(result)
			os.Stdout = stdOut

			// Assert
			tc.assertion(t, err, mockApi, output)
		})
	}
}
