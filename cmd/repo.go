package cmd

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

// Return cli command about repositories.
func (c *Cmd) RepoCommand() *cli.Command {
	return &cli.Command{
		Name:        "repo",
		Aliases:     []string{"r"},
		Description: "Get all repositories",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
		},
		Action: c.getRepositories,
	}
}

// Get all personal repositories.
// 1. If the github access token is set to Config,
//	 the target is all repositories (including private repos).
// 2. If the github access token is NOT set to Config,
//	 the target is public repositories (specify username as a "name" flag).
func (c *Cmd) getRepositories(cc *cli.Context) error {
	// With Github access token
	token := c.Config.Token
	if token != "" {
		rs, err := c.Api.ListRepositoriesForAuthenticatedUser()
		if err != nil {
			return err
		}
		for _, r := range rs {
			fmt.Println(r)
		}
		return nil
	}

	// With username
	useName := cc.String("name")
	if useName != "" {
		rs, err := c.Api.ListPublicRepositories(useName)
		if err != nil {
			return err
		}
		for _, r := range rs {
			fmt.Println(r)
		}
		return nil
	}

	// not correct usage
	return errors.New("Token or userName is not given.")
}
