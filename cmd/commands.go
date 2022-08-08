package cmd

import (
	"github.com/kokoichi206/go-git-stats/api"
	"github.com/kokoichi206/go-git-stats/util"
	"github.com/urfave/cli/v2"
)

// struct that has information related to subcommands
type Cmd struct {
	Config util.Config
	Api    api.ApiCaller
}

// Get all commands.
func (c *Cmd) NewCommands() []*cli.Command {
	return []*cli.Command{
		c.RepoCommand(),
		c.StatsCommand(),
		c.LinesCommand(),
	}
}
