package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const version = "0.0.0"

var (
	revision = "HEAD"
)

func main() {
	app := newApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newApp() *cli.App {
	return &cli.App{
		Name:    "ggs",
		Usage:   "Go git stats cli",
		Version: fmt.Sprintf("%s (rev:%s)", version, revision),
	}
}
