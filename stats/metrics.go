package stats

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Metrics map[string]int

type metric struct {
	Path  string
	Value int
}

func collectOpenedPullRequests(repository Repository, out chan<- metric) {
	pullRequests, err := repository.PullRequests("open", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- metric{Path: "pull_requests.open", Value: len(pullRequests)}
	if len(pullRequests) > 0 {
		value := int(time.Since(pullRequests[0].UpdatedAt).Hours() / 24)
		out <- metric{Path: "pull_requests.least_recently_updated_days", Value: value}
	}
}

func collectClosedPullRequests(repository Repository, out chan<- metric) {
	pullRequests, err := repository.PullRequests("closed", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- metric{Path: "pull_requests.closed", Value: len(pullRequests)}
}

func collectOpenedIssues(repository Repository, out chan<- metric) {
	issues, err := repository.Issues("open", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- metric{Path: "issues.open", Value: len(issues)}
}

func collectClosedIssues(repository Repository, out chan<- metric) {
	issues, err := repository.Issues("closed", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- metric{Path: "issues.closed", Value: len(issues)}
}

func (metrics Metrics) Compute(repository Repository) {
	feed := make(chan metric, 100)
	tasks := []func(Repository, chan<- metric){
		collectOpenedIssues,
		collectClosedIssues,
		collectOpenedPullRequests,
		collectClosedPullRequests,
	}

	var waitGrp sync.WaitGroup
	waitGrp.Add(len(tasks))
	for _, fn := range tasks {
		go func(fn func(Repository, chan<- metric)) {
			defer waitGrp.Done()
			fn(repository, feed)
		}(fn)
	}
	waitGrp.Wait()

	for {
		select {
		case m := <-feed:
			metrics[m.Path] = m.Value
		default:
			close(feed)
			return
		}
	}
}
