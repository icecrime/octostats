package metrics

import (
	"sync"
	"time"

	"github.com/icecrime/octostats/log"
	"github.com/icecrime/octostats/repository"
	"github.com/octokit/go-octokit/octokit"
)

type Metrics struct {
	Origin repository.Repository
	Items  []Metric
	m      sync.Mutex
}

func (m *Metrics) Add(items ...Metric) {
	m.m.Lock()
	defer m.m.Unlock()
	m.Items = append(m.Items, items...)
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

func collectIssues(issues []octokit.Issue) []Metric {
	var items []Metric
	for _, i := range issues {
		// Collect only issues that are not associated to pull requests.
		// All pull requests are issues but not all issues are pull requests.
		if i.PullRequest.HTMLURL == "" {
			m := NewMetric("issues.data", map[string]interface{}{
				"time":  i.CreatedAt.Unix(),
				"state": i.State,
				"id":    i.Number,
			})
			items = append(items, m)
		}
		for _, l := range i.Labels {
			m := NewMetric("labels.data", map[string]interface{}{
				"name": l.Name,
			})
			items = append(items, m)
		}
	}
	return items
}

func collectPrs(pullRequests []octokit.PullRequest) []Metric {
	var items []Metric
	for _, pr := range pullRequests {
		m := NewMetric("pull_requests.data", map[string]interface{}{
			"time":   pr.CreatedAt.Unix(),
			"state":  pr.State,
			"merged": pr.Merged,
			"id":     pr.ID,
		})
		items = append(items, m)
	}
	return items
}

func collectOpenedPullRequests(r repository.Repository) []Metric {
	pullRequests, err := r.PullRequests("open", "updated")
	if err != nil {
		log.Logger.Fatalf(err.Error())
	}

	var items []Metric

	items = append(items, NewMetric("pull_requests.open", map[string]interface{}{"count": len(pullRequests)}))
	if len(pullRequests) > 0 {
		items = append(items, collectPrs(pullRequests)...)

		value := int(time.Since(pullRequests[0].UpdatedAt).Hours() / 24)
		items = append(items, NewMetric("pull_requests.least_recently_updated_days", map[string]interface{}{"count": value}))
	}

	return items
}

func collectClosedPullRequests(r repository.Repository) []Metric {
	pullRequests, err := r.PullRequests("closed", "updated")
	if err != nil {
		log.Logger.Fatalf(err.Error())
	}
	var items []Metric
	items = append(items, NewMetric("pull_requests.closed", map[string]interface{}{"count": len(pullRequests)}))
	items = append(items, collectPrs(pullRequests)...)
	return items
}

func collectOpenedIssues(r repository.Repository) []Metric {
	issues, err := r.Issues("open", "updated")
	if err != nil {
		log.Logger.Fatalf(err.Error())
	}
	var items []Metric
	items = append(items, NewMetric("issues.open", map[string]interface{}{"count": len(issues)}))
	items = append(items, collectIssues(issues)...)
	return items
}

func collectClosedIssues(r repository.Repository) []Metric {
	issues, err := r.Issues("closed", "updated")
	if err != nil {
		log.Logger.Fatalf(err.Error())
	}
	var items []Metric
	items = append(items, NewMetric("issues.closed", map[string]interface{}{"count": len(issues)}))
	items = append(items, collectIssues(issues)...)
	return items
}

func Retrieve(r repository.Repository) *Metrics {
	tasks := []func(repository.Repository) []Metric{
		collectOpenedIssues,
		collectClosedIssues,
		collectOpenedPullRequests,
		collectClosedPullRequests,
	}

	var waitGrp sync.WaitGroup
	waitGrp.Add(len(tasks))

	metrics := New(r)

	for _, fn := range tasks {
		go func(fn func(repository.Repository) []Metric) {
			defer waitGrp.Done()
			metrics.Add(fn(r)...)
		}(fn)
	}
	waitGrp.Wait()

	log.Logger.Debug("Retrieve: end")
	return metrics
}
