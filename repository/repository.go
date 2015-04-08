package repository

import "github.com/octokit/go-octokit/octokit"

type Repository interface {
	Nwo() string
	Issues(string, string) ([]octokit.Issue, error)
	PullRequests(string, string) ([]octokit.PullRequest, error)
}
