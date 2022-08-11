package api

import "github.com/kokoichi206/go-git-stats/util"

func ExportNewApi(config util.Config) *Api {
	return New(config).(*Api)
}

func (a *Api) ExportGetConfig() util.Config {
	return a.config
}

func (a *Api) ExportSetURL(url string) {
	a.config.ApiBaseURL = url
}
