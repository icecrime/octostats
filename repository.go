package main

import (
	gh "github.com/crosbymichael/octokat"
)

type Repository interface {
	Id() gh.Repo
	Issues(string, string) ([]*gh.Issue, error)
	PullRequests(string, string) ([]*gh.PullRequest, error)
}
