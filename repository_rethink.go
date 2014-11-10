package main

import (
	"github.com/codegangsta/cli"
	gh "github.com/crosbymichael/octokat"
	rethink "github.com/dancannon/gorethink"
)

type rethinkRepository struct {
	id      *gh.Repo
	session *rethink.Session
}

func NewRethinkRepository(cli *cli.Context, repo *gh.Repo) (Repository, error) {
	return nil, nil
}

func (repository *rethinkRepository) Id() gh.Repo {
	return *repository.id
}

func (repository *rethinkRepository) Issues(status, sort string) ([]*gh.Issue, error) {
	return nil, nil
}

func (repository *rethinkRepository) PullRequests(status, sort string) ([]*gh.PullRequest, error) {
	return nil, nil
}
