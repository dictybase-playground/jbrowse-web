package jbrowse_manager

import (
	"net/http"
	"time"

	gh "github.com/google/go-github/v84/github"
)

type Config struct {
	Client *gh.Client
	Owner  string
	Repo   string
}

func NewConfig() Config {
	return Config{
		Client: gh.NewClient(&http.Client{Timeout: time.Second * 10}),
		Owner:  "GMOD",
		Repo:   "jbrowse-components",
	}
}
