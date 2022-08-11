package cmd

import (
	"sync"

	"github.com/kokoichi206/go-git-stats/api/mock"
	"github.com/kokoichi206/go-git-stats/util"
)

func ExportNewCommandWithMock(config util.Config, mockApi *mock.MockApi) Cmd {
	return Cmd{
		config: config,
		api:    mockApi,
		wait:   &sync.WaitGroup{},
		mutex:  &sync.Mutex{},
	}
}

func (c *Cmd) ExportInit() {
	c.total = 0
	c.config.Token = ""
}

func (c *Cmd) ExportSetToken(token string) {
	c.config.Token = token
}

func (c *Cmd) ExportGetTotal() int {
	return c.total
}
