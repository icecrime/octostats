package github

import (
	"testing"

	"github.com/icecrime/octostats/fixtures"
)

func TestAllPullRequests(t *testing.T) {
	fixtures.Setup()
	fixtures.SetupMux(t, "pulls")
	defer fixtures.TearDown()

	r := NewGitHubRepositoryWithClient("docker", "docker", fixtures.Client)
	prs, err := r.PullRequests("open", "updated")
	if err != nil {
		t.Fatal(err)
	}

	if len(prs) != 4 {
		t.Fatalf("Expected 4 prs but it was %d\n", len(prs))
	}
}

func TestAllIssues(t *testing.T) {
	fixtures.Setup()
	fixtures.SetupMux(t, "issues")
	defer fixtures.TearDown()

	r := NewGitHubRepositoryWithClient("docker", "docker", fixtures.Client)
	issues, err := r.Issues("open", "updated")
	if err != nil {
		t.Fatal(err)
	}

	if len(issues) != 4 {
		t.Fatalf("Expected 4 issues but it was %d\n", len(issues))
	}
}
