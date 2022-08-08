package main

import (
	"fmt"
	"log"
	"os"

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

	config := util.LoadConfig()
	api := api.New(config)

	cmd := cmd.Cmd{
		Config: config,
		Api:    api,
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
