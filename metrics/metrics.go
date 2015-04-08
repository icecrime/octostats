package metrics

import (
	"sync"
	"time"

	"github.com/icecrime/octostats/repository"

	log "github.com/Sirupsen/logrus"
)

type Metrics struct {
	Origin repository.Repository
	Items  map[string]int
}

func New(origin repository.Repository) *Metrics {
	return &Metrics{
		Origin: origin,
		Items:  make(map[string]int),
	}
}

type metric struct {
	Path  string
	Value int
}

func collectOpenedPullRequests(r repository.Repository, out chan<- metric) {
	pullRequests, err := r.PullRequests("open", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- metric{Path: "pull_requests.open", Value: len(pullRequests)}
	if len(pullRequests) > 0 {
		value := int(time.Since(pullRequests[0].UpdatedAt).Hours() / 24)
		out <- metric{Path: "pull_requests.least_recently_updated_days", Value: value}
	}
}

func collectClosedPullRequests(r repository.Repository, out chan<- metric) {
	pullRequests, err := r.PullRequests("closed", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- metric{Path: "pull_requests.closed", Value: len(pullRequests)}
}

func collectOpenedIssues(r repository.Repository, out chan<- metric) {
	issues, err := r.Issues("open", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- metric{Path: "issues.open", Value: len(issues)}
}

func collectClosedIssues(r repository.Repository, out chan<- metric) {
	issues, err := r.Issues("closed", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- metric{Path: "issues.closed", Value: len(issues)}
}

func Retrieve(r repository.Repository) *Metrics {
	feed := make(chan metric, 100)
	tasks := []func(repository.Repository, chan<- metric){
		collectOpenedIssues,
		collectClosedIssues,
		collectOpenedPullRequests,
		collectClosedPullRequests,
	}

	var waitGrp sync.WaitGroup
	waitGrp.Add(len(tasks))
	for _, fn := range tasks {
		go func(fn func(repository.Repository, chan<- metric)) {
			defer waitGrp.Done()
			fn(r, feed)
		}(fn)
	}
	waitGrp.Wait()

	metrics := New(r)
	for {
		select {
		case m := <-feed:
			metrics.Items[m.Path] = m.Value
		default:
			close(feed)
			return metrics
		}
	}
	return metrics
}
