package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

// Return cli command about repositories.
func (c *Cmd) StatsCommand() *cli.Command {
	return &cli.Command{
		Name:        "stats",
		Aliases:     []string{"s"},
		Description: "Get stats of a specific repository",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
		},
		Action: c.getStatistics,
	}
}

// Get statistics of a specific repository.
// If the github access token is set to Config,
// you can get stats of a private repository.
func (c *Cmd) getStatistics(cc *cli.Context) error {
	// get fullName (<userName>/<repo>)
	fullName := cc.String("name")
	if fullName == "" {
		// not correct usage
		return errors.New("name flag is not given.")
	}

	rs, err := c.api.WeeklyCommitActivity(fullName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("%-30s\t%-10s\t%-5s\n", "Start Time", "Additions", "Deletions")
	for _, r := range rs {
		fmt.Printf("%s\t%10d\t%5d\n", time.Unix(int64(r.Time), 0), r.Additions, r.Deletions)
	}
	return nil
}
