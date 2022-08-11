package cmd

import (
	"sync"

	"github.com/kokoichi206/go-git-stats/api"
	"github.com/kokoichi206/go-git-stats/util"
	"github.com/urfave/cli/v2"
)

// struct that has information related to subcommands
type Cmd struct {
	config util.Config
	api    api.ApiCaller
	wait   *sync.WaitGroup
	mutex  *sync.Mutex
	total  int
}

func New(config util.Config, api api.ApiCaller) Cmd {

	return Cmd{
		config: config,
		api:    api,
		wait:   &sync.WaitGroup{},
		mutex:  &sync.Mutex{},
		total:  0,
	}
}

// Get all commands.
func (c *Cmd) NewCommands() []*cli.Command {
	return []*cli.Command{
		c.RepoCommand(),
		c.StatsCommand(),
		c.LinesCommand(),
	}
}
