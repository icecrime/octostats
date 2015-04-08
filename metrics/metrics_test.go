package metrics

import (
	"testing"

	"github.com/icecrime/octostats/fixtures"
	"github.com/icecrime/octostats/github"
)

func TestCollectIssues(t *testing.T) {
	fixtures.Setup()
	fixtures.SetupMux(t, "issues")
	defer fixtures.TearDown()

	r := github.NewGitHubRepositoryWithClient("docker", "docker", fixtures.Client)
	m := New(r)

	feed := make(chan Metric, 100)

	collectOpenedIssues(r, feed)
	close(feed)

	for e := range feed {
		m.Items = append(m.Items, e)
	}

	// 1 global counter + 4 issues + 4 labels
	if len(m.Items) != 9 {
		t.Fatalf("Expected 8 metrics but got %d\n", len(m.Items))
	}
}
