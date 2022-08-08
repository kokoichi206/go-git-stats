package api

import (
	"github.com/kokoichi206/go-git-stats/util"
)

// struct that implements ApiCaller
type Api struct {
	config util.Config
}

func New(config util.Config) ApiCaller {
	return &Api{
		config: config,
	}
}
