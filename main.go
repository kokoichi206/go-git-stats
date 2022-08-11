package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kokoichi206/go-git-stats/api"
	"github.com/kokoichi206/go-git-stats/cmd"
	"github.com/kokoichi206/go-git-stats/util"
	"github.com/urfave/cli/v2"
)

const version = "0.0.0"

var (
	revision = "HEAD"
)

func main() {

	config, err := util.LoadConfig()
	if err != nil {
		fmt.Println(err.Error())
		// When the configuration error occures, the following steps should not be executed.
		os.Exit(1)
	}

	api := api.New(config)

	wg := sync.WaitGroup{}
	m := sync.Mutex{}

	cmd := cmd.Cmd{
		Config: config,
		Api:    api,
		Wait:   &wg,
		Mutex:  &m,
	}

	app := newApp(cmd)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newApp(c cmd.Cmd) *cli.App {
	return &cli.App{
		Name:     "ggs",
		Usage:    "Go git stats cli",
		Version:  fmt.Sprintf("%s (rev:%s)", version, revision),
		Commands: c.NewCommands(),
	}
}
