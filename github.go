package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/icecrime/octostats/repository"
	"github.com/octokit/go-octokit/octokit"
)

const (
	rateLimitRemaining = "X-RateLimit-Remaining"
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

func parseRepository(repo string) (string, string, error) {
	if splitRepos := strings.Split(repo, "/"); len(splitRepos) == 2 {
		return splitRepos[0], splitRepos[1], nil
	}
	return "", "", fmt.Errorf("bad repo format %s (expected username/repo)", repo)
}

func NewGitHubRepository(config *GitHubConfig) repository.Repository {
	token := githubAuthToken(config)
	ghClient := octokit.NewClient(&octokit.TokenAuth{AccessToken: token})

	owner, name, err := parseRepository(config.Repository)
	if err != nil {
		logger.Fatal(err)
	}
	return &gitHubRepository{Owner: owner, Name: name, client: ghClient}
}

type gitHubRepository struct {
	Owner  string
	Name   string
	client *octokit.Client
}

func (g *gitHubRepository) Nwo() string {
	return fmt.Sprintf("%s.%s", g.Owner, g.Name)
}

type IssuesCollection struct {
	Issues []octokit.Issue
	m      sync.Mutex
}

func (c *IssuesCollection) Add(issues ...octokit.Issue) {
	c.m.Lock()
	defer c.m.Unlock()

	c.Issues = append(c.Issues, issues...)
}

type PullRequestsCollection struct {
	PullRequests []octokit.PullRequest
	m            sync.Mutex
}

func (c *PullRequestsCollection) Add(prs ...octokit.PullRequest) {
	c.m.Lock()
	defer c.m.Unlock()

	c.PullRequests = append(c.PullRequests, prs...)
}

func (repo *gitHubRepository) Issues(state, sort string) ([]octokit.Issue, error) {
	u, err := repo.expandURL(octokit.RepoIssuesURL, state, sort)
	if err != nil {
		return nil, err
	}

	is := repo.client.Issues(u)
	first, res := is.All()
	coll := &IssuesCollection{Issues: first, m: sync.Mutex{}}

	if !res.HasError() && res.LastPage != nil {
		lastPage := res.LastPage
		u, _ := lastPage.Expand(nil)
		total, _ := strconv.Atoi(u.Query().Get("page"))

		if getRateLimitRemaining(res.Response) <= total {
			return coll.Issues, res.Err
		}

		urls := parseRemainingURLs(u, total)

		collectResults(urls, func(nu *url.URL) {
			is := repo.client.Issues(nu)
			next, res := is.All()

			if res.HasError() {
				logger.Debugf("Error fetching issues with %v\n", nu)
				return
			}
			coll.Add(next...)
		})

	}

	logger.Debugf("Loaded %d %s issues", len(coll.Issues), state)
	return coll.Issues, res.Err
}

func (repo *gitHubRepository) expandURL(link octokit.Hyperlink, state, sort string) (*url.URL, error) {
	queryParams := map[string]string{
		"sort":      sort,
		"direction": "asc",
		"state":     state,
		"per_page":  "100",
	}

	u, err := link.Expand(octokit.M{"owner": repo.Owner, "name": repo.Name})
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u, nil
}

func (repo *gitHubRepository) PullRequests(state, sort string) ([]octokit.PullRequest, error) {
	u, err := repo.expandURL(octokit.PullRequestsURL, state, sort)
	if err != nil {
		return nil, err
	}

	is := repo.client.PullRequests(u)
	first, res := is.All()
	coll := &PullRequestsCollection{PullRequests: first, m: sync.Mutex{}}

	if !res.HasError() && res.LastPage != nil {
		lastPage := res.LastPage
		u, _ := lastPage.Expand(nil)
		total, _ := strconv.Atoi(u.Query().Get("page"))

		if getRateLimitRemaining(res.Response) <= total {
			return coll.PullRequests, res.Err
		}

		urls := parseRemainingURLs(u, total)

		collectResults(urls, func(nu *url.URL) {
			is := repo.client.PullRequests(nu)
			next, res := is.All()

			if res.HasError() {
				logger.Debugf("Error fetching pull requests with %v\n", nu)
				return
			}
			coll.Add(next...)
		})

	}

	logger.Debugf("Loaded %d %s pull requests", len(coll.PullRequests), state)
	return coll.PullRequests, res.Err
}

func parseRemainingURLs(origin *url.URL, total int) []*url.URL {
	urls := make([]*url.URL, total-1)

	for i := 2; i <= total; i++ {
		np, _ := url.Parse(origin.String())

		q := np.Query()
		q.Set("page", strconv.Itoa(i))
		np.RawQuery = q.Encode()

		urls[i-2] = np
	}

	return urls
}

func collectResults(urls []*url.URL, collector func(*url.URL)) {
	var wg sync.WaitGroup
	for _, p := range urls {
		wg.Add(1)

		go func(nu *url.URL) {
			defer wg.Done()
			collector(nu)
		}(p)
	}
	wg.Wait()
}

func getRateLimitRemaining(res *octokit.Response) int {
	rate, err := strconv.Atoi(res.Header.Get(rateLimitRemaining))
	if err != nil {
		rate = 60
	}
	return rate
}
