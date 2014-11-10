package main

import (
	"io/ioutil"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	gh "github.com/crosbymichael/octokat"
)

func githubAuthToken(cli *cli.Context) string {
	token := cli.String("token")
	if cli.String("token-file") != "" {
		if fileContent, err := ioutil.ReadFile(cli.String("token-file")); err != nil {
			log.Fatalf(err.Error())
		} else {
			return string(fileContent)
		}
	}
	return token
}

func NewGitHubRepository(cli *cli.Context, id *gh.Repo) (Repository, error) {
	ghClient := gh.NewClient()
	ghClient.Token = githubAuthToken(cli)
	return &gitHubRepository{id: id, client: ghClient}, nil
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
		log.WithFields(log.Fields{"state": state, "page": page}).Debug("Loading issues")
		o.QueryParams["page"] = strconv.Itoa(page)
		if issues, err := repo.client.Issues(*repo.id, o); err != nil {
			return nil, err
		} else {
			prevSize = len(allIssues)
			allIssues = append(allIssues, issues...)
		}
	}
	log.Debugf("Loaded %d %s issues", len(allIssues), state)
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
		log.WithFields(log.Fields{"state": state, "page": page}).Debug("Loading pull requests")
		o.QueryParams["page"] = strconv.Itoa(page)
		if prs, err := repo.client.PullRequests(*repo.id, o); err != nil {
			return nil, err
		} else {
			prevSize = len(allPullRequests)
			allPullRequests = append(allPullRequests, prs...)
		}
	}
	log.Debugf("Loaded %d %s pull requests", len(allPullRequests), state)
	return allPullRequests, nil
}
