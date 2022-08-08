package cmd

import (
	"fmt"

	"github.com/kokoichi206/go-git-stats/api"
	"github.com/urfave/cli/v2"
)

// Return cli command about lines.
func (c *Cmd) LinesCommand() *cli.Command {
	return &cli.Command{
		Name:        "lines",
		Aliases:     []string{"l"},
		Description: "Get lines of codes you write before",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
		},
		Action: c.getLinesOfCodes,
	}
}

// Get lines of codes you write before.
func (c *Cmd) getLinesOfCodes(cc *cli.Context) error {

	var rs []api.Repository
	var err error

	// With Github access token
	token := c.Config.Token
	if token != "" {
		rs, err = c.Api.ListRepositoriesForAuthenticatedUser()
		if err != nil {
			return err
		}
	}

	// With username
	userName := cc.String("name")
	if userName != "" {
		rs, err = c.Api.ListPublicRepositories(userName)
		fmt.Println(rs)
		if err != nil {
			return err
		}
	}

	total := 0
	for _, r := range rs {
		st, err := c.Api.WeeklyCommitActivity(r.FullName)
		if err == nil {
			// Calculate the total amount of lines.
			for _, s := range st {
				total += s.Additions + s.Deletions
			}
		}
	}
	fmt.Println(total)

	return nil
}
