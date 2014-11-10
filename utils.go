package main

import (
	"fmt"
	"strings"

	gh "github.com/crosbymichael/octokat"
)

func parseRepository(repo string) (*gh.Repo, error) {
	if splitRepos := strings.Split(repo, "/"); len(splitRepos) == 2 {
		return &gh.Repo{Name: splitRepos[1], UserName: splitRepos[0]}, nil
	}
	return nil, fmt.Errorf("bad repo format %s (expected username/repo)", repo)
}
