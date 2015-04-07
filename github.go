package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/icecrime/octostats/repository"

	gh "github.com/crosbymichael/octokat"
)

func githubAuthToken(config *GitHubConfig) string {
	if config.AuthToken != "" {
		return config.AuthToken
	}

	fileContent, err := ioutil.ReadFile(config.AuthTokenFile)
	if err != nil {
		logger.WithField("error", err).WithField("filename", config.AuthTokenFile).Fatal("failed to load github auth token file")
	}
	return string(fileContent)
}

func logProgress(item, state string, page int) {
	logger.WithField("page", page).WithField("page", page).Debugf("loading %s", item)
}

func parseRepository(repo string) (*gh.Repo, error) {
	if splitRepos := strings.Split(repo, "/"); len(splitRepos) == 2 {
		return &gh.Repo{Name: splitRepos[1], UserName: splitRepos[0]}, nil
	}
	return nil, fmt.Errorf("bad repo format %s (expected username/repo)", repo)
}

func NewGitHubRepository(config *GitHubConfig) repository.Repository {
	ghClient := gh.NewClient()
	ghClient.Token = githubAuthToken(config)

	repoId, err := parseRepository(config.Repository)
	if err != nil {
		logger.Fatal(err)
	}
	return &gitHubRepository{id: repoId, client: ghClient}
}

type gitHubRepository struct {
	id     *gh.Repo
	client *gh.Client
}

func (repo *gitHubRepository) Id() gh.Repo {
	return *repo.id
}

func (repo *gitHubRepository) Issues(state, sort string) ([]*gh.Issue, error) {
	o := &gh.Options{}
	o.QueryParams = map[string]string{
		"sort":      sort,
		"direction": "asc",
		"state":     state,
		"per_page":  "100",
	}

	prevSize := -1
	allIssues := []*gh.Issue{}
	for page := 1; len(allIssues) != prevSize; page++ {
		logProgress("issues", state, page)
		o.QueryParams["page"] = strconv.Itoa(page)
		if issues, err := repo.client.Issues(*repo.id, o); err != nil {
			return nil, err
		} else {
			prevSize = len(allIssues)
			allIssues = append(allIssues, issues...)
		}
	}
	logger.Debugf("Loaded %d %s issues", len(allIssues), state)
	return allIssues, nil
}

func (repo *gitHubRepository) PullRequests(state, sort string) ([]*gh.PullRequest, error) {
	o := &gh.Options{}
	o.QueryParams = map[string]string{
		"sort":      sort,
		"direction": "asc",
		"state":     state,
		"per_page":  "100",
	}

	prevSize := -1
	allPullRequests := []*gh.PullRequest{}
	for page := 1; len(allPullRequests) != prevSize; page++ {
		logProgress("pull-requests", state, page)
		o.QueryParams["page"] = strconv.Itoa(page)
		if prs, err := repo.client.PullRequests(*repo.id, o); err != nil {
			return nil, err
		} else {
			prevSize = len(allPullRequests)
			allPullRequests = append(allPullRequests, prs...)
		}
	}
	logger.Debugf("Loaded %d %s pull requests", len(allPullRequests), state)
	return allPullRequests, nil
}
