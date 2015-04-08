package metrics

import (
	"sync"
	"time"

	"github.com/icecrime/octostats/repository"
	"github.com/octokit/go-octokit/octokit"

	log "github.com/Sirupsen/logrus"
)

type Metrics struct {
	Origin repository.Repository
	Items  []Metric
}

func New(origin repository.Repository) *Metrics {
	return &Metrics{
		Origin: origin,
		Items:  make([]Metric, 0),
	}
}

type Metric struct {
	Path string
	Data map[string]interface{}
}

func NewMetric(path string, data map[string]interface{}) Metric {
	return Metric{Path: path, Data: data}
}

func collectIssues(issues []octokit.Issue, out chan<- Metric) {
	for _, i := range issues {
		// Collect only issues that are not associated to pull requests.
		// All pull requests are issues but not all issues are pull requests.
		if i.PullRequest.HTMLURL == "" {
			m := NewMetric("issues.data", map[string]interface{}{
				"time":  i.CreatedAt.Unix(),
				"state": i.State,
				"id":    i.Number,
			})
			out <- m
		}
		for _, l := range i.Labels {
			m := NewMetric("labels.data", map[string]interface{}{
				"name": l.Name,
			})
			out <- m
		}
	}
}

func collectPrs(pullRequests []octokit.PullRequest, out chan<- Metric) {
	for _, pr := range pullRequests {
		m := NewMetric("pull_requests.data", map[string]interface{}{
			"time":   pr.CreatedAt.Unix(),
			"state":  pr.State,
			"merged": pr.Merged,
			"id":     pr.ID,
		})
		out <- m
	}
}

func collectOpenedPullRequests(r repository.Repository, out chan<- Metric) {
	pullRequests, err := r.PullRequests("open", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}

	out <- NewMetric("pull_requests.open", map[string]interface{}{"count": len(pullRequests)})
	if len(pullRequests) > 0 {
		collectPrs(pullRequests, out)

		value := int(time.Since(pullRequests[0].UpdatedAt).Hours() / 24)
		out <- NewMetric("pull_requests.least_recently_updated_days", map[string]interface{}{"count": value})
	}
}

func collectClosedPullRequests(r repository.Repository, out chan<- Metric) {
	pullRequests, err := r.PullRequests("closed", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- NewMetric("pull_requests.closed", map[string]interface{}{"count": len(pullRequests)})
	collectPrs(pullRequests, out)
}

func collectOpenedIssues(r repository.Repository, out chan<- Metric) {
	issues, err := r.Issues("open", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- NewMetric("issues.open", map[string]interface{}{"count": len(issues)})
	collectIssues(issues, out)
}

func collectClosedIssues(r repository.Repository, out chan<- Metric) {
	issues, err := r.Issues("closed", "updated")
	if err != nil {
		log.Fatalf(err.Error())
	}
	out <- NewMetric("issues.closed", map[string]interface{}{"count": len(issues)})
	collectIssues(issues, out)
}

func Retrieve(r repository.Repository) *Metrics {
	feed := make(chan Metric, 100)
	tasks := []func(repository.Repository, chan<- Metric){
		collectOpenedIssues,
		collectClosedIssues,
		collectOpenedPullRequests,
		collectClosedPullRequests,
	}

	var waitGrp sync.WaitGroup
	waitGrp.Add(len(tasks))
	for _, fn := range tasks {
		go func(fn func(repository.Repository, chan<- Metric)) {
			defer waitGrp.Done()
			fn(r, feed)
		}(fn)
	}
	waitGrp.Wait()

	metrics := New(r)
	for {
		select {
		case m := <-feed:
			metrics.Items = append(metrics.Items, m)
		default:
			close(feed)
			return metrics
		}
	}
	return metrics
}
