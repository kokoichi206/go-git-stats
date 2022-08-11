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

	var repositories []api.Repository
	var err error

	// With Github access token
	token := c.config.Token
	if token != "" {
		repositories, err = c.api.ListRepositoriesForAuthenticatedUser()
		if err != nil {
			return err
		}
	}

	// With username
	userName := cc.String("name")
	if userName != "" {
		repositories, err = c.api.ListPublicRepositories(userName)
		if err != nil {
			return err
		}
	}

	c.wait.Add(len(repositories))
	for _, repository := range repositories {
		go c.WeeklyCommitActivityAsyncCall(repository.FullName)
	}
	c.wait.Wait()

	// Final output
	fmt.Println(c.total)
	return nil
}

// Asynchronous API (WeeklyCommitActivity) call and calculate the total lines of codes.
func (c *Cmd) WeeklyCommitActivityAsyncCall(fullName string) {

	// Always decrements the WaitGroup counter.
	defer c.wait.Done()

	// Call function
	stats, err := c.api.WeeklyCommitActivity(fullName)
	if err != nil {
		return
	}

	// Calculate lines of codes of a specific repository.
	total := 0
	for _, s := range stats {
		total += s.Additions + s.Deletions
	}

	// Add to total lines of codes.
	c.mutex.Lock()
	c.total += total
	c.mutex.Unlock()

	return
}
